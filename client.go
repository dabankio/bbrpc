// refer to: https://github.com/btcsuite/btcd/blob/master/rpcclient/infrastructure.go

package bbrpc

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// sendPostBufferSize is the number of elements the HTTP POST send
	// channel can queue before blocking.
	sendPostBufferSize = 100
)

var (
	// ErrInvalidAuth is an error to describe the condition where the client
	// is either unable to authenticate or the specified endpoint is
	// incorrect.
	ErrInvalidAuth = errors.New("authentication failure")

	// ErrClientShutdown is an error to describe the condition where the
	// client is either already shutdown, or in the process of shutting
	// down.  Any outstanding futures when a client shutdown occurs will
	// return this error as will any new requests.
	ErrClientShutdown = errors.New("the client has been shutdown")
)

// response is the raw bytes of a JSON-RPC result, or the error if the response
// error object was non-null.
type response struct {
	result []byte
	err    error
}

type rawResponse struct {
	Result json.RawMessage `json:"result"`
	Error  error           `json:"error"`
}

// jsonRequest holds information about a json request that is used to properly
// detect, interpret, and deliver a reply to it.
type jsonRequest struct {
	id     uint64
	method string
	// cmd            interface{}
	marshalledJSON []byte
	responseChan   chan *response
}

// sendPostDetails houses an HTTP POST request to send to an RPC server as well
// as the original JSON-RPC command and a channel to reply on when the server
// responds with the result.
type sendPostDetails struct {
	httpRequest *http.Request
	jsonRequest *jsonRequest
}

// ConnConfig describes the connection configuration parameters for the client.
// This
type ConnConfig struct {
	// Host is the IP address and port of the RPC server you want to connect
	// to.
	Host string

	// User is the username to use to authenticate to the RPC server.
	User string

	// Pass is the passphrase to use to authenticate to the RPC server.
	Pass string

	// DisableTLS specifies whether transport layer security should be
	// disabled.  It is recommended to always use TLS if the RPC server
	// supports it as otherwise your username and password is sent across
	// the wire in cleartext.
	DisableTLS bool

	// Certificates are the bytes for a PEM-encoded certificate chain used
	// for the TLS connection.  It has no effect if the DisableTLS parameter
	// is true.
	Certificates []byte
}

// Client .
type Client struct {
	Debug bool //print some debug log when true

	id uint64 // atomic, so must stay 64-bit aligned
	// config holds the connection configuration associated with this client.
	config     *ConnConfig
	httpClient *http.Client

	// Track command and their response channels by ID.
	requestLock sync.Mutex

	sendPostChan chan *sendPostDetails
	shutdown     chan struct{}

	wg sync.WaitGroup
}

// NextID returns the next id to be used when sending a JSON-RPC message.  This
// ID allows responses to be associated with particular requests per the
// JSON-RPC specification.  Typically the consumer of the client does not need
// to call this function, however, if a custom request is being created and used
// this function should be used to ensure the ID is unique amongst all requests
// being made.
func (c *Client) NextID() uint64 {
	return atomic.AddUint64(&c.id, 1)
}

// sendPost sends the passed request to the server by issuing an HTTP POST
// request using the provided response channel for the reply.  Typically a new
// connection is opened and closed for each command when using this method,
// however, the underlying HTTP client might coalesce multiple commands
// depending on several factors including the remote server configuration.
func (c *Client) sendPost(jReq *jsonRequest) {
	// Generate a request to the configured RPC server.
	url := c.config.Host
	if !strings.HasPrefix(c.config.Host, "http") {
		protocol := "http"
		if !c.config.DisableTLS {
			protocol = "https"
		}
		url = protocol + "://" + c.config.Host
	}
	bodyReader := bytes.NewReader(jReq.marshalledJSON)
	httpReq, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		jReq.responseChan <- &response{result: nil, err: err}
		return
	}
	// httpReq.Close = true
	httpReq.Header.Set("Content-Type", "application/json")

	if c.config.User != "" {
		// Configure basic access authorization.
		httpReq.SetBasicAuth(c.config.User, c.config.Pass)
	}

	// log.Printf("Sending command [%s] with id %d", jReq.method, jReq.id)

	select {
	case <-c.shutdown:
		jReq.responseChan <- &response{result: nil, err: ErrClientShutdown}
	default:
	}

	c.sendPostChan <- &sendPostDetails{
		jsonRequest: jReq,
		httpRequest: httpReq,
	}
}

// newHTTPClient returns a new http client that is configured according to the
// proxy and TLS settings in the associated connection configuration.
func newHTTPClient(config *ConnConfig) (*http.Client, error) {
	// Set proxy function if there is a proxy configured.
	// var proxyFunc func(*http.Request) (*url.URL, error)
	// Configure TLS if needed.
	var tlsConfig *tls.Config
	if !config.DisableTLS {
		if len(config.Certificates) > 0 {
			pool := x509.NewCertPool()
			pool.AppendCertsFromPEM(config.Certificates)
			tlsConfig = &tls.Config{
				RootCAs: pool,
			}
		}
	}

	client := http.Client{
		Timeout: time.Second * 15,
		Transport: &http.Transport{
			// Proxy:           proxyFunc,
			TLSClientConfig: tlsConfig,
		},
	}

	return &client, nil
}

// sendPostHandler handles all outgoing messages when the client is running
// in HTTP POST mode.  It uses a buffered channel to serialize output messages
// while allowing the sender to continue running asynchronously.  It must be run
// as a goroutine.
func (c *Client) sendPostHandler() {
out:
	for {
		// Send any messages ready for send until the shutdown channel
		// is closed.
		select {
		case details := <-c.sendPostChan:
			c.handleSendPostMessage(details)
		case <-c.shutdown:
			break out
		}
	}

	// Drain any wait channels before exiting so nothing is left waiting
	// around to send.
cleanup:
	for {
		select {
		case details := <-c.sendPostChan:
			details.jsonRequest.responseChan <- &response{
				result: nil,
				err:    ErrClientShutdown,
			}

		default:
			break cleanup
		}
	}
	c.wg.Done()
	if c.Debug {
		log.Printf("RPC client send handler done for %s", c.config.Host)
	}
}

// handleSendPostMessage handles performing the passed HTTP request, reading the
// result, unmarshalling it, and delivering the unmarshalled result to the
// provided response channel.
func (c *Client) handleSendPostMessage(details *sendPostDetails) {
	jReq := details.jsonRequest
	// log.Printf("Sending post [%s] with id %d, json: %s", jReq.method, jReq.id, string(jReq.marshalledJSON))
	httpResponse, err := c.httpClient.Do(details.httpRequest)
	if httpResponse != nil {
		defer httpResponse.Body.Close()
	}
	if err != nil {
		jReq.responseChan <- &response{err: err}
		return
	}

	// Read the raw bytes and close the response.
	respBytes, err := ioutil.ReadAll(httpResponse.Body)
	// httpResponse.Body.Close()
	if err != nil {
		err = fmt.Errorf("error reading json reply: %v", err)
		jReq.responseChan <- &response{err: err}
		return
	}

	if c.Debug {
		log.Println("[bbrpc dbg] resp", string(respBytes))
	}
	// Try to unmarshal the response as a regular JSON-RPC response.
	var resp rawResponse
	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		// When the response itself isn't a valid JSON-RPC response
		// return an error which includes the HTTP status code and raw
		// response bytes.
		err = fmt.Errorf("status code: %d, response: %q", httpResponse.StatusCode, string(respBytes))
		jReq.responseChan <- &response{err: err}
		return
	}

	res, err := resp.Result.MarshalJSON()
	jReq.responseChan <- &response{result: res, err: err}
}

// NewClient .
func NewClient(config *ConnConfig) (*Client, error) {
	var httpClient *http.Client
	var err error

	httpClient, err = newHTTPClient(config)
	if err != nil {
		return nil, err
	}

	client := &Client{
		config:       config,
		httpClient:   httpClient,
		sendPostChan: make(chan *sendPostDetails, sendPostBufferSize),
		shutdown:     make(chan struct{}),
	}

	client.wg.Add(1)
	go client.sendPostHandler()
	return client, nil
}

// NewClientWith .
func NewClientWith(config *ConnConfig, httpClient *http.Client) (*Client, error) {
	client := &Client{
		config:       config,
		httpClient:   httpClient,
		sendPostChan: make(chan *sendPostDetails, sendPostBufferSize),
		shutdown:     make(chan struct{}),
	}

	client.wg.Add(1)
	go client.sendPostHandler()
	return client, nil
}

// Shutdown shuts down the client by disconnecting any connections associated
// with the client and, when automatic reconnect is enabled, preventing future
// attempts to reconnect.  It also stops all goroutines.
func (c *Client) Shutdown() {
	// Do the shutdown under the request lock to prevent clients from
	// adding new requests while the client shutdown process is initiated.
	c.requestLock.Lock()
	defer c.requestLock.Unlock()

	// Ignore the shutdown request if the client is already in the process
	// of shutting down or already shutdown.
	if !c.doShutdown() {
		return
	}
}

// doShutdown closes the shutdown channel and logs the shutdown unless shutdown
// is already in progress.  It will return false if the shutdown is not needed.
//
// This function is safe for concurrent access.
func (c *Client) doShutdown() bool {
	// Ignore the shutdown request if the client is already in the process
	// of shutting down or already shutdown.
	select {
	case <-c.shutdown:
		return false
	default:
	}

	if c.Debug {
		log.Printf("Shutting down RPC client %s", c.config.Host)
	}
	close(c.shutdown)
	c.httpClient.CloseIdleConnections()
	return true
}

// futureParse receives from the passed futureResult channel to extract a
// reply or any errors.  The examined errors include an error in the
// futureResult and the error in the reply from the server.  This will block
// until the result is available on the passed channel.
func futureParse(f chan *response, v interface{}) error {
	// Wait for a response on the returned channel.
	r := <-f
	if r.err != nil {
		return fmt.Errorf("RPC request return error: %v", r.err)
	}
	err := json.Unmarshal(r.result, v)
	if err != nil {
		err = fmt.Errorf("failed to parse rpc json response %s to %t, %v", string(r.result), v, err)
	}
	return err
}

// CallJSONRPC any rpc
func (c *Client) CallJSONRPC(method string, params map[string]interface{}, result interface{}) ([]byte, error) {
	resp, err := c.sendCmd(method, params)
	if err != nil {
		return nil, err
	}
	if result != nil {
		return nil, futureParse(resp, &result)
	}
	r := <-resp
	return r.result, r.err
}

// sendCmd sends the passed command to the associated server and returns a
// response channel on which the reply will be delivered at some point in the
// future.  It handles both websocket and HTTP POST mode depending on the
// configuration of the client.
func (c *Client) sendCmd(method string, param interface{}) (chan *response, error) {
	id := c.NextID()

	var rawParams json.RawMessage
	if param == nil {
		rawParams = json.RawMessage("{}")
	} else {
		marshalledParam, err := json.Marshal(param)
		if err != nil {
			return nil, err
		}
		rawParams = json.RawMessage(marshalledParam)
	}

	req := &Request{
		Jsonrpc: "2.0",
		ID:      id,
		Method:  method,
		Params:  rawParams,
	}

	if c.Debug {
		log.Println("[bbrpc dbg] req:", method, string(req.Params))
	}

	marshalledJSON, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Generate the request and send it along with a channel to respond on.
	responseChan := make(chan *response, 1)
	jReq := &jsonRequest{
		id:             id,
		method:         method,
		marshalledJSON: marshalledJSON,
		responseChan:   responseChan,
	}
	c.sendPost(jReq)
	return responseChan, nil
}

// PendingPostCount current post request count waiting for process
func (c *Client) PendingPostCount() int {
	return len(c.sendPostChan)
}
