package bbrpc

// Getblockcount https://github.com/bigbangcore/BigBang/wiki/JSON-RPC#getblockcount
func (c *Client) Getblockcount(fork *string) (*int64, error) {
	resp, err := c.sendCmd("getblockcount", struct {
		Fork *string `json:"fork,omitempty"`
	}{Fork: fork})
	if err != nil {
		return nil, err
	}
	var height int64
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
func (c *Client) Getforkheight(fork *string) (h int64, err error) {
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
