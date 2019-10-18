package bbrpc

import (
	"encoding/json"
	"fmt"
	"log"
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

// TestBlockTime 该测试用于观察出块周期,运行n分钟，然后打印区块信息
func TestBlockTime(t *testing.T) {
	killBigBangServer, client, _ := tRunBigbangServerAndBeginMint(t)
	defer killBigBangServer()

	time.Sleep(time.Minute * 6) //基本确保矿工有至少15*3个币

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

	fmt.Println(toJSONIndent(blocks))
}

// 计算出块间隔
// 结论：第一个块通常约39s,后续每35s一个块，可以用41s作为测试用的大致出块间隔
func TestBlockTimeDuration(t *testing.T) {
	_json := `
	[
		{
		  "hash": "3b6f275b4fe20c7945395755a1cf93a96e17dc64edf752e1c7cdbef14c008e0d",
		  "version": 1,
		  "type": "primary-pow",
		  "time": 1571033632,
		  "fork": "69ed2d6ebb9ace57ca0c1e6eaa4d8eaa2db458bfaa37b37d1a9bfcec46ef099e",
		  "height": 1,
		  "txmint": "5da41220ec7d54637f98577865e0338e5f8ded3a5c763d5695e6e9162ebc1cf2",
		  "prev": "69ed2d6ebb9ace57ca0c1e6eaa4d8eaa2db458bfaa37b37d1a9bfcec46ef099e",
		  "tx": []
		},
		{
		  "hash": "81b7da9cf884007d85218d7a8f00e6e80fb9ec294f962e5dafb3d28363dab051",
		  "version": 1,
		  "type": "primary-pow",
		  "time": 1571033671,
		  "fork": "69ed2d6ebb9ace57ca0c1e6eaa4d8eaa2db458bfaa37b37d1a9bfcec46ef099e",
		  "height": 2,
		  "txmint": "5da41247632bcc98c304d7f27af08e3beb309028dec748deebcd4b3e154dc0ef",
		  "prev": "3b6f275b4fe20c7945395755a1cf93a96e17dc64edf752e1c7cdbef14c008e0d",
		  "tx": []
		},
		{
		  "hash": "8eecde857e52ff67f399122b62ab38235d8b4c8ca2f6e343ad4a4a8cf3d6247f",
		  "version": 1,
		  "type": "primary-pow",
		  "time": 1571033706,
		  "fork": "69ed2d6ebb9ace57ca0c1e6eaa4d8eaa2db458bfaa37b37d1a9bfcec46ef099e",
		  "height": 3,
		  "txmint": "5da4126a9a7d418f8c7a730d97aa99228ff280039b5d738a04f87f27c704e9c6",
		  "prev": "81b7da9cf884007d85218d7a8f00e6e80fb9ec294f962e5dafb3d28363dab051",
		  "tx": []
		},
		{
		  "hash": "9843e0459f37a6f0713f630735171096071bbcda5ada661ccb2f9281667d6172",
		  "version": 1,
		  "type": "primary-pow",
		  "time": 1571033741,
		  "fork": "69ed2d6ebb9ace57ca0c1e6eaa4d8eaa2db458bfaa37b37d1a9bfcec46ef099e",
		  "height": 4,
		  "txmint": "5da4128d8f49101ccd807fc3be13b3ef3d3fa2ffd65d706e132303b5f7a3d630",
		  "prev": "8eecde857e52ff67f399122b62ab38235d8b4c8ca2f6e343ad4a4a8cf3d6247f",
		  "tx": []
		},
		{
		  "hash": "12ab23082616c1e00a1e6f8524864495522e729fb2103736ee400aa0bdd4eb4b",
		  "version": 1,
		  "type": "primary-pow",
		  "time": 1571033776,
		  "fork": "69ed2d6ebb9ace57ca0c1e6eaa4d8eaa2db458bfaa37b37d1a9bfcec46ef099e",
		  "height": 5,
		  "txmint": "5da412b0a206dda0bf8a08e5c510e916a1aacab3d7deaeb3d7e8e894d9724d33",
		  "prev": "9843e0459f37a6f0713f630735171096071bbcda5ada661ccb2f9281667d6172",
		  "tx": []
		},
		{
		  "hash": "66a8519d3633c062685e6c3c7319ab6782ffc864e96eb298fd341d31f068e6e2",
		  "version": 1,
		  "type": "primary-pow",
		  "time": 1571033811,
		  "fork": "69ed2d6ebb9ace57ca0c1e6eaa4d8eaa2db458bfaa37b37d1a9bfcec46ef099e",
		  "height": 6,
		  "txmint": "5da412d38bb539dbb134559cfe1aeca77aecca931508e5804902f600913f0442",
		  "prev": "12ab23082616c1e00a1e6f8524864495522e729fb2103736ee400aa0bdd4eb4b",
		  "tx": []
		},
		{
		  "hash": "9ca9154070e2ae03d86e6e0f20ceb2cd7d35a3c83a546b9dea23892e54e4dee9",
		  "version": 1,
		  "type": "primary-pow",
		  "time": 1571033846,
		  "fork": "69ed2d6ebb9ace57ca0c1e6eaa4d8eaa2db458bfaa37b37d1a9bfcec46ef099e",
		  "height": 7,
		  "txmint": "5da412f65bd00c62e37c91d6e5af1a8d315eb7a64bbd76cbeb868d6cc605e111",
		  "prev": "66a8519d3633c062685e6c3c7319ab6782ffc864e96eb298fd341d31f068e6e2",
		  "tx": []
		},
		{
		  "hash": "9cb09dfe9b4df58adacc208640a0abbe144166a7662aa1f61020abcfc2c61b62",
		  "version": 1,
		  "type": "primary-pow",
		  "time": 1571033881,
		  "fork": "69ed2d6ebb9ace57ca0c1e6eaa4d8eaa2db458bfaa37b37d1a9bfcec46ef099e",
		  "height": 8,
		  "txmint": "5da41319ee1822744e1280cab47c55fee1cf3f5b2310961d5b1aa8bd267fdc06",
		  "prev": "9ca9154070e2ae03d86e6e0f20ceb2cd7d35a3c83a546b9dea23892e54e4dee9",
		  "tx": []
		},
		{
		  "hash": "2488454428d5231c2c690ca765b5205e58912de6392c7bd9bbe363073d3b6763",
		  "version": 1,
		  "type": "primary-pow",
		  "time": 1571033916,
		  "fork": "69ed2d6ebb9ace57ca0c1e6eaa4d8eaa2db458bfaa37b37d1a9bfcec46ef099e",
		  "height": 9,
		  "txmint": "5da4133cac9051cdc7db12842dfccff50adebe8bbbe9a7fffdedb2e5f2eb4677",
		  "prev": "9cb09dfe9b4df58adacc208640a0abbe144166a7662aa1f61020abcfc2c61b62",
		  "tx": []
		},
		{
		  "hash": "264a6e2d584f308e785af6f6066b4a94f4df08052d23d414fabca57c54d95427",
		  "version": 1,
		  "type": "primary-pow",
		  "time": 1571033951,
		  "fork": "69ed2d6ebb9ace57ca0c1e6eaa4d8eaa2db458bfaa37b37d1a9bfcec46ef099e",
		  "height": 10,
		  "txmint": "5da4135f7e1ce988e07386d048d4a186455abd6ea49373fc56f93dccb404b161",
		  "prev": "2488454428d5231c2c690ca765b5205e58912de6392c7bd9bbe363073d3b6763",
		  "tx": []
		},
		{
		  "hash": "190dcfea28b018a4bb2cd5aeb85e0719edf95adbee6d74fc3aeb69d69d30b0ab",
		  "version": 1,
		  "type": "primary-pow",
		  "time": 1571033986,
		  "fork": "69ed2d6ebb9ace57ca0c1e6eaa4d8eaa2db458bfaa37b37d1a9bfcec46ef099e",
		  "height": 11,
		  "txmint": "5da4138232cfb18ec6a708d61ec67d5539bd43bf56cdc216cf2844b42c32ee9d",
		  "prev": "264a6e2d584f308e785af6f6066b4a94f4df08052d23d414fabca57c54d95427",
		  "tx": []
		}
	  ]
	`
	var blocks []BlockInfo

	err := json.Unmarshal([]byte(_json), &blocks)
	tShouldNil(t, err)

	for i, block := range blocks {
		t := time.Unix(int64(block.Time), 0)
		var du time.Duration
		if i > 0 {
			lastTime := time.Unix(int64(blocks[i-1].Time), 0)
			du = t.Sub(lastTime)
		}
		fmt.Println(t, du)
	}
}

// 测试代币的发行/挖矿/交易/查询/遍历
func TestToken(t *testing.T) {

}
