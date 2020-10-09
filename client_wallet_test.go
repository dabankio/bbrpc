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

		count := len(listk)

		_, err = c.Importprivkey("514025fb4b6d6bdb15d4521d047d20ace5311fa10e2e8889adbd262f93dc673b", "123", nil)
		tShouldNil(t, err)

		listk, err = c.Listkey()
		tShouldNil(t, err)
		tShouldTrue(t, len(listk) == count+1)
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

func TestClient_Exportkey(t *testing.T) {
	testClientMethod(t, func(c *Client) {
		p, err := c.Makekeypair()
		tShouldNil(t, err)

		fmt.Printf("pair: %#v\n", p)

		_, err = c.Importprivkey(p.Privkey, "123", nil)
		tShouldNil(t, err)

		key, err := c.Exportkey(p.Pubkey)
		tShouldNil(t, err)
		fmt.Println("export key:", key)
	})
}

func TestClient_Importkey(t *testing.T) {
	testClientMethod(t, func(c *Client) {
		// pair: &bbrpc.Keypair{Privkey:"52b18397cfe2464d80df629ea35050a34f51379d1a38f39a94f8838f2b731b66", Pubkey:"56005784cf72f3a4e228e7dd4459f3c00e69ea731d8d4b07327d4d619fc37909"}
		// export key: 0979c39f614d7d32074b8d1d73ea690ec0f35944dde728e2a4f372cf8457005601000000aaca6c71d4f72cbad55ae8a698f813deb2f8d6706db8114a3c6f7553790142f869c79867b5c1fb33ea0946747400c9f90c271b1f35061026cd2c82c0

		pubkey := "0979c39f614d7d32074b8d1d73ea690ec0f35944dde728e2a4f372cf8457005601000000aaca6c71d4f72cbad55ae8a698f813deb2f8d6706db8114a3c6f7553790142f869c79867b5c1fb33ea0946747400c9f90c271b1f35061026cd2c82c0"
		key, err := c.Importkey(pubkey, nil)
		tShouldNil(t, err)
		fmt.Println("importkey ", key, err)

		keys, err := c.Listkey()
		tShouldNil(t, err)
		fmt.Printf("listkeys %#v\n", keys)

		ret, err := c.Unlockkey(keys[0].Key, "123", nil)
		tShouldNil(t, err)
		fmt.Println("unlock ret:", *ret)
	})
}

func TestClient_Addnewtemplate(t *testing.T) {
	tw := TW{T: t}
	tw.Continue(true)
	runClientTest(t, func(c *Client, minerAddr string) {
		pubks := make([]string, 2)
		for i := 0; i < len(pubks); i++ {
			p, err := c.Makekeypair()
			tw.Nil(err)
			pubks[i] = p.Pubkey
		}

		tplAddr, err := c.Addnewtemplate(AddnewtemplateParamMultisig{
			Required: 2, Pubkeys: pubks,
		})
		tw.Nil(err)
		fmt.Println(tplAddr)
	})
}
