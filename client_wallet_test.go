package bbrpc

import (
	"encoding/json"
	"fmt"
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

func TestClient_Validateaddress(t *testing.T) {
	testClientMethod(t, func(c *Client) {
		ret, err := c.Validateaddress("1ah1herg25ggkgv73ewdwaxjdbzd2r2mtrzbk8pvfsxzbg7srf8rrtqnc")
		tShouldNil(t, err)
		fmt.Println("add is valid?", toJSONIndent(ret))
	})
}

func TestJSONUnmarshal(t *testing.T) {
	data := `{"isvalid":true,"addressdata":{"address":"1gav1f0tqtdy5j6jeayybrv4nyffmk2bmwt74wdyh04vthv3khc9qmsgg","ismine":true,"type":"pubkey","pubkey":"138b73eca83701d1374e8ee6748949dff3956cbcbc574e1a597cd3578317b682"}}`
	var info AddressInfo
	err := json.Unmarshal([]byte(data), &info)
	tShouldNil(t, err)
	fmt.Println(toJSONIndent(info))
}
