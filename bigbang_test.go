package bbrpc

import (
	"fmt"
	"log"
	"testing"
	"time"
)

// 测试pow挖矿,简单挖矿并列出余额
func TestPOWMine(t *testing.T) {
	killBigBangServer, client, templateAddress := tRunBigbangServerAndBeginMint(t)
	defer killBigBangServer()

	listk, err := client.Listkey()
	tShouldNil(t, err)
	tShouldTrue(t, len(listk) == 2, listk)

	tWait4mine()

	balance, err := client.Getbalance(nil, nil)
	tShouldNil(t, err)

	fmt.Println("balance:", toJSONIndent(balance))
	fmt.Println("addr:", toJSONIndent(tCryptonightAddr))
	fmt.Println("key:", toJSONIndent(tCryptonightKey))

	{ //尝试把挖到的币花费掉
		// result, err := client.Unlockkey(tCryptonightKey.Pubkey, _tPassphrase, nil)
		// tShouldNil(t, err)
		// tShouldTrue(t, strings.Contains(*result, "success"))

		// result, err = client.Unlockkey(tCryptonightAddr.Pubkey, _tPassphrase, nil)
		// tShouldNil(t, err)
		// tShouldTrue(t, strings.Contains(*result, "success"))

		txid, err := client.Sendfrom(CmdSendfrom{
			From:   templateAddress,
			To:     tCryptonightAddr.Address,
			Amount: 50,
		})
		tShouldNil(t, err)
		fmt.Println("sendfrom txid", *txid)
	}
	tWait4mine()

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

// 挖矿，不停的打印blockcount和balance
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

	time.Sleep(time.Second * 10) //一段时间后停止
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

	var addrs []AddrKeypair
	for i := 0; i < 4; i++ {
		a := makeKeyPairAddr(client, t)
		addrs = append(addrs, a)
	}
	for _, a := range addrs {
		fmt.Printf("%#v\n", a)
	}
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

	tWait4mine()
	tWait4mine()

	for _, k := range []AddrKeypair{tAddr0, tAddr1} {
		ret, err := client.Importprivkey(k.Privkey, _tPassphrase)
		tShouldNil(t, err)
		tShouldTrue(t, ret != nil)
		tShouldNotContains(t, *ret, "error")

		_, err = client.Unlockkey(k.Pubkey, _tPassphrase, nil)
		tShouldNil(t, err)
	}

	balOfMiner, err := client.Getbalance(nil, &mintTplAddress)
	tShouldNil(t, err)
	tShouldTrue(t, balOfMiner[0].Avail > 45, balOfMiner)

	for i := 0; i < 3; i++ {
		txid, err := client.Sendfrom(CmdSendfrom{
			From: mintTplAddress, To: tAddr0.Address, Amount: 15,
		})
		tShouldNil(t, err)
		tShouldTrue(t, txid != nil)
	}

	tWait4mine()

	{ // 0 transfer to 1
		txid, err := client.Sendfrom(CmdSendfrom{
			From: tAddr0.Address, To: tAddr1.Address, Amount: 32,
		})
		tShouldNil(t, err)
		tShouldTrue(t, txid != nil)
	}

	tWait4mine()

	bal, err := client.Getbalance(nil, nil)
	tShouldNil(t, err)
	fmt.Println("balance", toJSONIndent(bal))

	hash1, err := client.Getblockhash(1, nil)
	tShouldNil(t, err)
	_, txs, err := client.ListTransactionsSinceBlock(hash1[0], 10)
	tShouldNil(t, err)
	fmt.Println("tx...", toJSONIndent(txs))

	//余额不足的情况
	_, err = client.Sendfrom(CmdSendfrom{
		From: tAddr0.Address, To: tAddr1.Address, Amount: 1000,
	})
	tShouldTrue(t, err != nil)
	fmt.Println("insufficient error", err)
}

// TestBlockTime 该测试用于观察出块周期,运行n分钟，然后打印区块信息
// 计算出块间隔
// 结论：基于1.0.0版本，大部分在3s内，少部分最多30s
func TestBlockTime(t *testing.T) {
	killBigBangServer, client, mintTplAddress := tRunBigbangServerAndBeginMint(t)
	defer killBigBangServer()

	time.Sleep(time.Minute * 1)

	balOfMiner, err := client.Getbalance(nil, &mintTplAddress)
	tShouldNil(t, err)
	tShouldTrue(t, balOfMiner[0].Avail > 45, balOfMiner)

	fk, err := client.Getforkheight(nil)
	tShouldNil(t, err)

	var blocks []*BlockInfo
	for h := 1; h <= int(fk); h++ {
		hash, err := client.Getblockhash(h, nil)
		tShouldNil(t, err)
		b, err := client.Getblock(hash[0])
		tShouldNil(t, err)
		blocks = append(blocks, b)
	}

	// fmt.Println(toJSONIndent(blocks))
	for i, block := range blocks {
		tm := time.Unix(int64(block.Time), 0)
		var du time.Duration
		if i > 0 {
			lastTime := time.Unix(int64(blocks[i-1].Time), 0)
			du = tm.Sub(lastTime)
		}
		fmt.Println(tm, du)
	}
}

// 测试代币的发行/挖矿/交易/查询/遍历
func TestToken(t *testing.T) {

}
