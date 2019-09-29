package bbrpc

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
