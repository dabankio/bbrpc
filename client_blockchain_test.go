package bbrpc

import (
	"fmt"
	"testing"
)

func TestClient_Listfork(t *testing.T) {
	testClientMethod(t, func(client *Client) {
		forks, err := client.Listfork(true)
		tShouldNil(t, err, "failed to listfork")
		tShouldTrue(t, len(forks) > 0, "zero forks len")
		fmt.Printf("forks: \n%#v\n", forks)
	})
}

func TestClient_Getforkheight(t *testing.T) {
	testClientMethod(t, func(client *Client) {
		height, err := client.Getforkheight(nil)
		tShouldNil(t, err)
		tShouldTrue(t, height >= 0, "bad height", height)
	})
}

func TestClient_Getblockcount(t *testing.T) {
	testClientMethod(t, func(c *Client) {
		count, err := c.Getblockcount(nil)
		tShouldNil(t, err)
		tShouldTrue(t, count != nil)
		fmt.Println("block count:", *count)
	})
}

func TestClient_Getblock_Getblockhash_Getforkheight_Listfork(t *testing.T) {
	testClientMethod(t, func(c *Client) {
		Wait4nBlocks(1, c)

		h, err := c.Getforkheight(nil)
		tShouldNil(t, err)

		ha, err := c.Getblockhash(int(h), nil)
		tShouldNil(t, err)

		b, err := c.Getblock(ha[0])
		tShouldNil(t, err)
		fmt.Println("block", toJSONIndent(b))

		forks, err := c.Listfork(true)
		tShouldNil(t, err)
		fmt.Println("forks", toJSONIndent(forks))
	})
}
