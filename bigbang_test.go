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
	cryptonightAddr := AddrKeypair{}
	cryptonightKey := AddrKeypair{}
	passphrase := "123"

	prepareAddress := func() {
		defer fmt.Println("------------prepare mining data finished------------")
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

		cryptonightAddr = makeKeyPairAddr(client, t)
		cryptonightKey = makeKeyPairAddr(client, t)

		tShouldNotZero(t, cryptonightAddr.Address)
		tShouldNotZero(t, cryptonightAddr.Privkey)
		tShouldNotZero(t, cryptonightAddr.Pubkey)

		tShouldNotZero(t, cryptonightKey.Privkey)
		tShouldNotZero(t, cryptonightKey.Pubkey)
	}

	prepareAddress()

	// time.Sleep(time.Second * 2) //debug port usage

	runBBOptions := DefaultDebugBBArgs()
	runBBOptions["cryptonightaddress"] = &cryptonightAddr.Address
	runBBOptions["cryptonightkey"] = &cryptonightKey.Privkey

	killBigBangServer, err := RunBigBangServer(&RunBigBangOptions{
		NewTmpDir: true,
		Args:      runBBOptions,
	})
	tShouldNil(t, err, "failed to run bigbang server")
	defer killBigBangServer()

	client, err := NewClient(DefaultDebugConnConfig())
	tShouldNil(t, err, "failed to new rpc client")
	defer client.Shutdown()

	_, err = client.Importprivkey(cryptonightAddr.Privkey, passphrase)
	tShouldNil(t, err)
	_, err = client.Importprivkey(cryptonightKey.Privkey, passphrase)
	tShouldNil(t, err)

	listk, err := client.Listkey()
	tShouldNil(t, err)
	tShouldTrue(t, len(listk) == 2, listk)

	templateAddress, err := client.Addnewtemplate(AddnewtemplateParamMint{
		Mint:  cryptonightKey.Pubkey,
		Spent: cryptonightAddr.Address,
	})
	tShouldNil(t, err)
	fmt.Println("mint template address", *templateAddress)

	time.Sleep(time.Minute)
	// time.Sleep(time.Second * 15)

	balance, err := client.Getbalance(nil, nil)
	tShouldNil(t, err)

	fmt.Println("balance:", toJSONIndent(balance))
	fmt.Println("addr:", toJSONIndent(cryptonightAddr))
	fmt.Println("key:", toJSONIndent(cryptonightKey))

	{ //尝试把挖到的币花费掉
		result, err := client.Unlockkey(cryptonightKey.Pubkey, passphrase, nil)
		tShouldNil(t, err)
		tShouldTrue(t, strings.Contains(*result, "success"))

		result, err = client.Unlockkey(cryptonightAddr.Pubkey, passphrase, nil)
		tShouldNil(t, err)
		tShouldTrue(t, strings.Contains(*result, "success"))

		// txdata, err := client.Createtransaction(CmdCreatetransaction{
		// 	From:   *templateAddress,
		// 	To:     cryptonightAddr.Address,
		// 	Amount: 23,
		// })
		// tShouldNil(t, err)
		// fmt.Println("txdata", *txdata)

		txid, err := client.Sendfrom(CmdSendfrom{
			From:   *templateAddress,
			To:     cryptonightAddr.Address,
			Amount: 0.005,
		})
		tShouldNil(t, err)
		fmt.Println("sendfrom txid", *txid)
	}
	time.Sleep(time.Minute) //休眠1分钟等待打包

	{
		// balance, err := client.Getbalance(nil, nil)
		// tShouldNil(t, err)

		// fmt.Println("balance:", toJSONIndent(balance))
	}

	forkHeight, err := client.Getforkheight(nil)
	tShouldNil(t, err)
	fmt.Println("fork height", forkHeight)

	{ // 尝试迭代整个链
		fmt.Println("--------尝试迭代整个链---------")
		printHeight := func(height int) {
			hash, err := client.Getblockhash(int64(height), nil)
			tShouldNil(t, err)
			tShouldTrue(t, len(hash) == 1)

			block, err := client.Getblock(hash[0])
			tShouldNil(t, err)

			fmt.Println("---------height:", height, "--------")
			fmt.Println(toJSONIndent(block))
		}
		for i := 0; i < int(forkHeight); i++ {
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

// 测试交易（挖矿产出）
func TestTransaction(t *testing.T) {

}

func TestRunMineNode(t *testing.T) {
	cryptonightAddr := AddrKeypair{}
	cryptonightKey := AddrKeypair{}

	prepareAddress := func() {
		defer fmt.Println("------------prepare mining data finished------------")
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

		cryptonightAddr = makeKeyPairAddr(client, t)
		cryptonightKey = makeKeyPairAddr(client, t)
	}

	prepareAddress()

	// time.Sleep(time.Second * 2) //debug port usage
	passphrase := "123"
	_ = passphrase

	runBBOptions := DefaultDebugBBArgs()
	runBBOptions["cryptonightaddress"] = &cryptonightAddr.Address
	runBBOptions["cryptonightkey"] = &cryptonightKey.Privkey

	killBigBangServer, err := RunBigBangServer(&RunBigBangOptions{
		NewTmpDir: true,
		Args:      runBBOptions,
	})
	tShouldNil(t, err, "failed to run bigbang server")
	defer killBigBangServer()

	client, err := NewClient(DefaultDebugConnConfig())
	tShouldNil(t, err, "failed to new rpc client")
	defer client.Shutdown()

	// _, err = client.Importprivkey(cryptonightAddr.Privkey, passphrase)
	// tShouldNil(t, err)
	// _, err = client.Importprivkey(cryptonightKey.Privkey, passphrase)
	// tShouldNil(t, err)

	templateAddress, err := client.Addnewtemplate(AddnewtemplateParamMint{
		Mint:  cryptonightKey.Pubkey,
		Spent: cryptonightAddr.Address,
	})
	tShouldNil(t, err)
	fmt.Println("mint template address", *templateAddress)

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

func TestPrepareMineAddress(t *testing.T) {
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
}

func makeKeyPairAddr(c *Client, t *testing.T) AddrKeypair {
	k, err := c.Makekeypair()
	tShouldNil(t, err)

	add, err := c.Getpubkeyaddress(k.Pubkey, nil)
	tShouldNil(t, err)

	return AddrKeypair{Keypair: *k, Address: *add}
}
