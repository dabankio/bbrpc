package bbrpc

import (
	"fmt"
	"testing"
)

func TestViewTestnet(t *testing.T) {
	t.SkipNow()

	client, err := NewClient(&ConnConfig{
		Host:       "127.0.0.1:9904",
		DisableTLS: true,
	})
	tShouldNil(t, err, "failed to new rpc client!")
	defer client.Shutdown()

	var (
		maxHeight int64
	)
	{
		maxHeight, err = client.Getforkheight(nil)
		tShouldNil(t, err)
		fmt.Println("forkHeight", maxHeight)
	}

	{ //listransaction
		txs, err := client.Listtransaction(nil, nil)
		tShouldNil(t, err)
		fmt.Println("list tx...", toJSONIndent(txs))
	}

	{ //recent 10 blocks
		for i := maxHeight; i > maxHeight-10; i-- {
			hash, err := client.Getblockhash(int(i), nil)
			tShouldNil(t, err)
			blk, err := client.Getblock(hash[0])
			tShouldNil(t, err)

			fmt.Println("----------------height", i)
			fmt.Println("block", toJSONIndent(blk))

			if len(blk.Tx) == 0 {
				fmt.Println("::no tx this block")
			} else {
				fmt.Println("tx...")
				for _, txid := range blk.Tx {
					tx, err := client.Gettransaction(txid, pbool(false))
					tShouldNil(t, err)
					fmt.Println(toJSONIndent(tx))
				}
			}

		}
	}
}
