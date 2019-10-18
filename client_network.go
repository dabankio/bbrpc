package bbrpc

// Listpeer https://github.com/bigbangcore/BigBang/wiki/JSON-RPC#listpeer
func (c *Client) Listpeer() ([]PeerInfo, error) {
	resp, err := c.sendCmd("listpeer", nil)
	if err != nil {
		return nil, err
	}
	var data []PeerInfo
	err = futureParse(resp, &data)
	return data, err
}
