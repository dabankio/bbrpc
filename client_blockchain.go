package bbrpc

// Getblockcount https://github.com/bigbangcore/BigBang/wiki/JSON-RPC#getblockcount
func (c *Client) Getblockcount(fork *string) (*int, error) {
	resp, err := c.sendCmd("getblockcount", struct {
		Fork *string `json:"fork,omitempty"`
	}{Fork: fork})
	if err != nil {
		return nil, err
	}
	var height int
	err = futureParse(resp, &height)
	return &height, err
}

// Getblock https://github.com/bigbangcore/BigBang/wiki/JSON-RPC#getblock
func (c *Client) Getblock(hash string) (*BlockInfo, error) {
	resp, err := c.sendCmd("getblock", struct {
		Block string `json:"block"`
	}{Block: hash})
	if err != nil {
		return nil, err
	}
	var info BlockInfo
	err = futureParse(resp, &info)
	return &info, err
}

// Getblockdetail https://github.com/bigbangcore/BigBang/wiki/JSON-RPC#getblockdetail
func (c *Client) Getblockdetail(blockHash string) (*BlockDetail, error) {
	resp, err := c.sendCmd("getblockdetail", struct {
		Block string `json:"block"`
	}{blockHash})
	if err != nil {
		return nil, err
	}
	var info BlockDetail
	err = futureParse(resp, &info)
	return &info, err
}

// Getblockhash https://github.com/bigbangcore/BigBang/wiki/JSON-RPC#getblockhash
func (c *Client) Getblockhash(height int, fork *string) ([]string, error) {
	resp, err := c.sendCmd("getblockhash", struct {
		Height int     `json:"height"`
		Fork   *string `json:"fork,omitempty"`
	}{Height: height, Fork: fork})
	if err != nil {
		return nil, err
	}
	var hash []string
	err = futureParse(resp, &hash)
	return hash, err
}

// Getforkheight https://github.com/bigbangcore/BigBang/wiki/JSON-RPC#getforkheight
func (c *Client) Getforkheight(fork *string) (h int, err error) {
	resp, err := c.sendCmd("getforkheight", struct {
		Fork *string `json:"fork,omitempty"`
	}{Fork: fork})
	if err != nil {
		return -1, err
	}
	err = futureParse(resp, &h)
	return
}

// Listfork https://github.com/bigbangcore/BigBang/wiki/JSON-RPC#listfork
func (c *Client) Listfork(all bool) (ret []ForkProfile, err error) {
	resp, err := c.sendCmd("listfork", struct {
		All bool `json:"all"`
	}{All: all})
	if err != nil {
		return nil, err
	}
	err = futureParse(resp, &ret)
	return
}

// Sendtransaction https://github.com/bigbangcore/BigBang/wiki/JSON-RPC#sendtransaction
func (c *Client) Sendtransaction(txdata string) (*string, error) {
	resp, err := c.sendCmd("sendtransaction", struct {
		Txdata string `json:"txdata"`
	}{Txdata: txdata})
	if err != nil {
		return nil, err
	}
	var txid string
	err = futureParse(resp, &txid)
	return &txid, err
}
