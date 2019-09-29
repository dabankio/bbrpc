package bbrpc

// Version https://github.com/bigbangcore/BigBang/wiki/JSON-RPC#version
func (c *Client) Version() (string, error) {
	resp, err := c.sendCmd("version", nil)
	if err != nil {
		return "", err
	}
	var ver string
	err = futureParse(resp, &ver)
	return ver, err
}
