package bbrpc

import (
	"testing"
)

func TestClient_Importprivkey(t *testing.T) {
	testClientMethod(t, func(c *Client) {
		listk, err := c.Listkey()
		tShouldNil(t, err)
		tShouldTrue(t, len(listk) == 0)

		_, err = c.Importprivkey("514025fb4b6d6bdb15d4521d047d20ace5311fa10e2e8889adbd262f93dc673b", "123")
		tShouldNil(t, err)

		listk, err = c.Listkey()
		tShouldNil(t, err)
		tShouldTrue(t, len(listk) == 1)
	})
}
