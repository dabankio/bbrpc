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

func TestClient_Getblock_Getblockhash_Getforkheight_Listfork_Getblockdetail(t *testing.T) {
	tw := TW{T: t}
	_pass := "123"
	runClientTest(t, func(c *Client, minerAddr string) {
		tw.Nil(Wait4balanceReach(minerAddr, 100, c))

		nk, err := c.Getnewkey(_pass)
		tw.Nil(err)
		add, err := c.Getpubkeyaddress(nk, nil)
		tw.Nil(err)

		txid, err := c.Sendfrom(CmdSendfrom{
			From:   minerAddr,
			To:     *add,
			Amount: 12,
		})
		tw.Nil(err)

		// tw.Nil(Wait4nBlocks(1, c))
		tw.Nil(Wait4balanceReach(*add, 10, c))

		h, err := c.Getforkheight(nil)
		tShouldNil(t, err)

		txs, err := c.Gettransaction(*txid, nil)
		tw.Nil(err)

		// fmt.Println("tx:", h, toJSONIndent(txs))
		ha, err := c.Getblockhash(int(h)+1-*txs.Transaction.Confirmations, nil)
		tShouldNil(t, err)

		b, err := c.Getblock(ha[0])
		tShouldNil(t, err)
		tw.AllNotZero("下列字段应该非空", b.Hash, b.Version, b.Type, b.Time, b.Fork, b.Height, b.Prev)

		blockDetail, err := c.Getblockdetail(ha[0])
		tw.Nil(err).NotZero(blockDetail)
		d := blockDetail
		tw.AllNotZero("下列字段应该不为空", d.Hash, d.HashPrev, d.Version, d.Type, d.Time, d.Bits, d.Fork, d.Height, d.Prev, d.Txmint)
		mt := d.Txmint
		tw.AllNotZero("下列字段应该不为空", mt.Txid, mt.Version, mt.Type, mt.Time, mt.Anchor, mt.Sendto, mt.Amount, mt.Fork)
		if len(d.Tx) == 0 { //有可能取到的块不包含这比交易，rpc间隙间又出了块
			fmt.Println("len(d.Tx) == 0 ??", toJSONIndent(d))
		}
		tx := d.Tx[0]
		tw.AllNotZero("下列字段应该不为空", tx.Txid, tx.Version, tx.Type, tx.Time, tx.Anchor, tx.Sendfrom, tx.Sendto, tx.Amount, tx.Txfee, tx.Sig, tx.Confirmations, tx.Fork)

		{
			forks, err := c.Listfork(true)
			tShouldNil(t, err)
			tw.True(len(forks) == 1)
			f := forks[0]
			tw.AllNotZero("fork fields not zero", f.Fork, f.Name, f.Reward, f.Owner)
		}
	})
}
