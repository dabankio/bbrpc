package bbrpc

import (
	"fmt"
	"testing"
)

func TestClient_Makekeypair(t *testing.T) {
	testClientMethod(t, func(c *Client) {
		pair, err := c.Makekeypair()
		tShouldNil(t, err)
		tShouldTrue(t, pair.Pubkey != "", "empty pubkey")
		fmt.Printf("%#v\n", pair)
	})
}

func TestClient_Getpubkeyaddress(t *testing.T) {
	pair := &Keypair{Privkey: "24d936e04b93c95ca29c5afa5a71cec24f4f5e5290621f8c0e4ebac4c7a095a9", Pubkey: "cd2eaef7ed048ea211b37d196dc63446840df5a796b571cda7a6cfb4d536c74d"}
	shouldBe := "19q3kdndmsykafkbhppbafx8dgh339hkd35yv64d2hr2evxxe5v6gc8jm"
	testClientMethod(t, func(c *Client) {
		add, err := c.Getpubkeyaddress(pair.Pubkey, nil)
		tShouldNil(t, err)
		tShouldTrue(t, *add == shouldBe)
		fmt.Println("addr", *add)
	})
}
