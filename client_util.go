package bbrpc

// Decodetransaction https://github.com/bigbangcore/BigBang/wiki/JSON-RPC#decodetransaction
func (c *Client) Decodetransaction(txdata string) (*NoneSerializedTransaction, error) {
	resp, err := c.sendCmd("decodetransaction", struct {
		Txdata string `json:"txdata"`
	}{Txdata: txdata})
	if err != nil {
		return nil, err
	}
	var ret NoneSerializedTransaction
	err = futureParse(resp, &ret)
	return &ret, err
}

// Getpubkeyaddress https://github.com/bigbangcore/BigBang/wiki/JSON-RPC#getpubkeyaddress
func (c *Client) Getpubkeyaddress(pubkey string, reversal *string) (*string, error) {
	resp, err := c.sendCmd("getpubkeyaddress", struct {
		Pubkey   string  `json:"pubkey,omitempty"`
		Reversal *string `json:"reversal,omitempty"`
	}{Pubkey: pubkey, Reversal: reversal})
	if err != nil {
		return nil, err
	}
	var data string
	err = futureParse(resp, &data)
	return &data, err
}

// Listunspent https://github.com/bigbangcore/BigBang/wiki/JSON-RPC#listunspent
func (c *Client) Listunspent(address string, fork *string, max uint) (*UnspentTotal, error) {
	resp, err := c.sendCmd("listunspent", struct {
		Address string  `json:"address"`
		Fork    *string `json:"fork,omitempty"`
		Max     uint    `json:"max"`
	}{
		Address: address, Fork: fork, Max: max,
	})
	if err != nil {
		return nil, err
	}
	var data UnspentTotal
	err = futureParse(resp, &data)
	return &data, err
}

// Makekeypair https://github.com/bigbangcore/BigBang/wiki/JSON-RPC#makekeypair
func (c *Client) Makekeypair() (*Keypair, error) {
	resp, err := c.sendCmd("makekeypair", nil)
	if err != nil {
		return nil, err
	}
	var data Keypair
	err = futureParse(resp, &data)
	return &data, err
}
