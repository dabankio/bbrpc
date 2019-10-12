package bbrpc

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
)

// 测试pow挖矿
func TestPOWMine(t *testing.T) {
	killBigBangServer, client, templateAddress := tRunBigbangServerAndBeginMint(t)
	defer killBigBangServer()

	listk, err := client.Listkey()
	tShouldNil(t, err)
	tShouldTrue(t, len(listk) == 2, listk)

	time.Sleep(time.Second * 25)

	balance, err := client.Getbalance(nil, nil)
	tShouldNil(t, err)

	fmt.Println("balance:", toJSONIndent(balance))
	fmt.Println("addr:", toJSONIndent(tCryptonightAddr))
	fmt.Println("key:", toJSONIndent(tCryptonightKey))

	{ //尝试把挖到的币花费掉
		result, err := client.Unlockkey(tCryptonightKey.Pubkey, _tPassphrase, nil)
		tShouldNil(t, err)
		tShouldTrue(t, strings.Contains(*result, "success"))

		result, err = client.Unlockkey(tCryptonightAddr.Pubkey, _tPassphrase, nil)
		tShouldNil(t, err)
		tShouldTrue(t, strings.Contains(*result, "success"))

		txid, err := client.Sendfrom(CmdSendfrom{
			From:   templateAddress,
			To:     tCryptonightAddr.Address,
			Amount: 50,
		})
		tShouldNil(t, err)
		fmt.Println("sendfrom txid", *txid)
	}
	time.Sleep(time.Second * 45) //休眠一段时间等待打包

	forkHeight, err := client.Getforkheight(nil)
	tShouldNil(t, err)
	fmt.Println("fork height", forkHeight)

	{ // 尝试迭代整个链
		fmt.Println("--------尝试迭代整个链---------")
		printHeight := func(height int) {
			hash, err := client.Getblockhash(height, nil)
			tShouldNil(t, err)
			tShouldTrue(t, len(hash) == 1)

			block, err := client.Getblock(hash[0])
			tShouldNil(t, err)

			fmt.Println("---------height:", height, "--------")
			fmt.Println(toJSONIndent(block))
			if len(block.Tx) > 0 {
				for _, txid := range block.Tx {
					tx, err := client.Gettransaction(txid, nil)
					tShouldNil(t, err)
					fmt.Println("tx", toJSONIndent(tx))
				}
			}
		}
		for i := 0; i <= int(forkHeight); i++ {
			printHeight(i)
		}
	}

	{ //尝试列出所有的交易
		txs, err := client.Listtransaction(nil, nil)
		tShouldNil(t, err)

		for _, tx := range txs {
			fmt.Println("tx", toJSONIndent(tx))

			txDetail, err := client.Gettransaction(tx.Txid, pbool(false))
			tShouldNil(t, err)
			fmt.Println("tx detail", toJSONIndent(txDetail))
		}
	}

	{
		balance, err := client.Getbalance(nil, nil)
		tShouldNil(t, err)

		fmt.Println("balance:", toJSONIndent(balance))
	}
}

func TestRunMineNode(t *testing.T) {
	killBigBangServer, client, _ := tRunBigbangServerAndBeginMint(t)
	defer killBigBangServer()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	tickerDone := make(chan bool)

	type gotBlock struct {
		Count int64
		Time  time.Time
	}
	gotBlocks := make([]gotBlock, 0)
	go func() {
		for {
		slct:
			select {
			case <-tickerDone:
				fmt.Println("Done!")
				return
			case tm := <-ticker.C:
				count, err := client.Getblockcount(nil)
				tShouldNil(t, err)

				for _, b := range gotBlocks {
					if b.Count == *count {
						break slct
					}
				}
				gotBlocks = append(gotBlocks, gotBlock{*count, tm})
				log.Println("[blockcount]", toJSONIndent(gotBlocks))

				bal, err := client.Getbalance(nil, nil)
				tShouldNil(t, err)
				log.Println("balance", toJSONIndent(bal))
			}
		}
	}()

	time.Sleep(time.Minute * 1)
	tickerDone <- true
}

//准备2组地址
func TestPrepare2Address(t *testing.T) {
	killBigBangServer, err := RunBigBangServer(&RunBigBangOptions{
		NewTmpDir:       true,
		NotPrint2stdout: true,
		Args:            DefaultDebugBBArgs(),
	})
	tShouldNil(t, err, "failed to run bigbang server")
	defer killBigBangServer()

	client, err := NewClient(DefaultDebugConnConfig())
	tShouldNil(t, err, "failed to new rpc client")
	defer client.Shutdown()

	a0 := makeKeyPairAddr(client, t)
	a1 := makeKeyPairAddr(client, t)

	fmt.Println(toJSONIndent(a0), toJSONIndent(a1))
	fmt.Printf("%#v\n%#v\n", a0, a1)
}

func makeKeyPairAddr(c *Client, t *testing.T) AddrKeypair {
	k, err := c.Makekeypair()
	tShouldNil(t, err)

	add, err := c.Getpubkeyaddress(k.Pubkey, nil)
	tShouldNil(t, err)

	return AddrKeypair{Keypair: *k, Address: *add}
}

// 测试单笔交易需要多个vin (交易额大于单个utxo的情况)
// 使用0-5地址
// 给0转入3资金，每笔15
// 0 transfer to 1, 32
func TestMultiVinTx(t *testing.T) {
	killBigBangServer, client, mintTplAddress := tRunBigbangServerAndBeginMint(t)
	defer killBigBangServer()

	time.Sleep(time.Minute) //基本确保矿工有至少15*3个币

	for _, k := range []AddrKeypair{tAddr0, tAddr1} {
		ret, err := client.Importprivkey(k.Privkey, _tPassphrase)
		tShouldNil(t, err)
		tShouldTrue(t, ret != nil)
		tShouldNotContains(t, *ret, "error")

		_, err = client.Unlockkey(k.Pubkey, _tPassphrase, nil)
		tShouldNil(t, err)
	}

	for i := 0; i < 3; i++ {
		txid, err := client.Sendfrom(CmdSendfrom{
			From: mintTplAddress, To: tAddr0.Address, Amount: 15,
		})
		tShouldNil(t, err)
		tShouldTrue(t, txid != nil)
	}

	time.Sleep(time.Minute) //等待打包

	{ // 0 transfer to 1
		txid, err := client.Sendfrom(CmdSendfrom{
			From: tAddr0.Address, To: tAddr1.Address, Amount: 32,
		})
		tShouldNil(t, err)
		tShouldTrue(t, txid != nil)
	}

	time.Sleep(time.Minute) //等待打包

	bal, err := client.Getbalance(nil, nil)
	tShouldNil(t, err)
	fmt.Println("balance", toJSONIndent(bal))

	hash1, err := client.Getblockhash(1, nil)
	tShouldNil(t, err)
	_, txs, err := client.ListTransactionsSinceBlock(hash1[0])
	tShouldNil(t, err)
	fmt.Println("tx...", toJSONIndent(txs))

	//余额不足的情况
	_, err = client.Sendfrom(CmdSendfrom{
		From: tAddr0.Address, To: tAddr1.Address, Amount: 1000,
	})
	tShouldTrue(t, err != nil)
	fmt.Println("insufficient error", err)
}
