package bbrpc

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestSimpleSendfromWithData(t *testing.T) {
	killBigBangServer, client, templateAddress := TesttoolRunServerAndBeginMint(t)
	defer killBigBangServer()

	tShouldNil(t, Wait4balanceReach(templateAddress, 100, client))

	a := TAddr0
	_, err := client.Importprivkey(a.Privkey, _tPassphrase)
	tShouldNil(t, err)
	_, err = client.Unlockkey(a.Pubkey, _tPassphrase, nil)
	tShouldNil(t, err)

	_, err = client.Sendfrom(CmdSendfrom{
		From:   templateAddress,
		To:     a.Address,
		Amount: 12.3,
		Data:   pstring("0xfab1"),
	})
	tShouldNil(t, err)

	tShouldNil(t, Wait4balanceReach(a.Address, 12, client))

	toAddr := TAddr1
	txid, err := client.Sendfrom(CmdSendfrom{
		From:   a.Address,
		To:     toAddr.Address,
		Amount: 12.1,
		Data:   pstring(UtilDataEncoding("卢本伟🐂🍺")),
	})
	tShouldNil(t, err)
	tShouldNil(t, Wait4nBlocks(1, client))

	// tx, err := client.Gettransaction(*txid, pbool(true))
	tx, err := client.Gettransaction(*txid, nil)
	tShouldNil(t, err)
	fmt.Println("tx::", toJSONIndent(tx))
	fmt.Println(UtilDataDecoding(tx.Transaction.Data))

	if tx.Serialization != nil {
		detx, err := client.Decodetransaction(*tx.Serialization)
		tShouldNil(t, err)
		fmt.Println("tx::", toJSONIndent(detx))
	}
	time.Sleep(2 * time.Second)

}
func TestSimpleSignTX(t *testing.T) {
	killBigBangServer, client, templateAddress := TesttoolRunServerAndBeginMint(t)
	defer killBigBangServer()

	tShouldNil(t, Wait4balanceReach(templateAddress, 100, client))

	a := TAddr0
	_, err := client.Importprivkey(a.Privkey, _tPassphrase)
	tShouldNil(t, err)
	_, err = client.Unlockkey(a.Pubkey, _tPassphrase, nil)
	tShouldNil(t, err)

	_, err = client.Sendfrom(CmdSendfrom{
		From:   templateAddress,
		To:     a.Address,
		Amount: 12.3,
	})
	tShouldNil(t, err)

	tShouldNil(t, Wait4balanceReach(a.Address, 12, client))

	ret, err := client.Createtransaction(CmdCreatetransaction{
		From:   a.Address,
		To:     TAddr1.Address,
		Amount: 12.1,
		Data:   pstring(UtilDataEncoding("abc123")),
	})
	tShouldNil(t, err)
	tShouldNotZero(t, ret)
	fmt.Println("created tx:", *ret)
	fmt.Println("privk:", a.Privkey)

	sret, err := client.Signtransaction(*ret)
	tShouldNil(t, err)
	tShouldNotZero(t, sret)

	txid, err := client.Sendtransaction(sret.Hex)
	tShouldNil(t, err)
	tShouldNil(t, Wait4nBlocks(1, client))

	tx, err := client.Gettransaction(*txid, nil)
	tShouldNil(t, err)
	fmt.Println("tx::", toJSONIndent(tx))

	// DecodeTx 这个，有问题（钱包那边），暂时就先不测试了
	// deRet, err := client.Decodetransaction(sret.Hex)
	// tShouldNil(t, err)
	// tShouldNotZero(t, deRet)
	// fmt.Println("decode tx:", toJSONIndent(*deRet))
	// fmt.Println("sig:", deRet.Sig)
	// fmt.Println("signed hex:", sret.Hex)
}

// 测试pow挖矿,简单挖矿并列出余额
func TestPOWMine(t *testing.T) {
	killBigBangServer, client, templateAddress := TesttoolRunServerAndBeginMint(t)
	defer killBigBangServer()

	listk, err := client.Listkey()
	tShouldNil(t, err)
	tShouldTrue(t, len(listk) == 2, listk)

	Wait4nBlocks(1, client)

	balance, err := client.Getbalance(nil, nil)
	tShouldNil(t, err)

	fmt.Println("balance:", toJSONIndent(balance))
	fmt.Println("addr:", toJSONIndent(TCryptonightAddr))
	fmt.Println("key:", toJSONIndent(TCryptonightKey))

	{ //尝试把挖到的币花费掉
		// result, err := client.Unlockkey(TCryptonightKey.Pubkey, _tPassphrase, nil)
		// tShouldNil(t, err)
		// tShouldTrue(t, strings.Contains(*result, "success"))

		// result, err = client.Unlockkey(TCryptonightAddr.Pubkey, _tPassphrase, nil)
		// tShouldNil(t, err)
		// tShouldTrue(t, strings.Contains(*result, "success"))

		txid, err := client.Sendfrom(CmdSendfrom{
			From:   templateAddress,
			To:     TCryptonightAddr.Address,
			Amount: 50,
		})
		tShouldNil(t, err)
		fmt.Println("sendfrom txid", *txid)
	}
	Wait4nBlocks(1, client)

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
	killBigBangServer, client, _ := TesttoolRunServerAndBeginMint(t)
	defer killBigBangServer()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	tickerDone := make(chan bool)

	type gotBlock struct {
		Count int
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
	killBigBangServer, client, mintTplAddress := TesttoolRunServerAndBeginMint(t)
	defer killBigBangServer()

	Wait4nBlocks(1, client)

	for _, k := range []AddrKeypair{TAddr0, TAddr1} {
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
			From: mintTplAddress, To: TAddr0.Address, Amount: 15,
		})
		tShouldNil(t, err)
		tShouldTrue(t, txid != nil)
	}

	Wait4nBlocks(1, client)

	{ // 0 transfer to 1
		txid, err := client.Sendfrom(CmdSendfrom{
			From: TAddr0.Address, To: TAddr1.Address, Amount: 32,
		})
		tShouldNil(t, err)
		tShouldTrue(t, txid != nil)
	}

	Wait4nBlocks(1, client)

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
		From: TAddr0.Address, To: TAddr1.Address, Amount: 1000,
	})
	tShouldTrue(t, err != nil)
	fmt.Println("insufficient error", err)
}

// TestBlockTime 该测试用于观察出块周期,运行n分钟，然后打印区块信息
// 计算出块间隔
// 结论：基于1.0.0版本，大部分在3s内，少部分最多30s
func TestBlockTime(t *testing.T) {
	killBigBangServer, client, mintTplAddress := TesttoolRunServerAndBeginMint(t)
	defer killBigBangServer()

	Wait4nBlocks(5, client)

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

// 测试单个节点2个地址的多重签名
func TestMultisigSingleNode(t *testing.T) {
	killBigBangServer, client, minerAddress := TesttoolRunServerAndBeginMint(t)
	defer killBigBangServer()

	tShouldNil(t, Wait4balanceReach(minerAddress, 100, client))

	// 使用2个地址，产生一个多签地址
	// 将资金转入多签地址
	// 从多签地址将资金转出
	a0, a1, a2 := TAddr0, TAddr1, TAddr2
	for _, a := range []AddrKeypair{a0, a1, a2} {
		_, err := client.Importprivkey(a.Privkey, _tPassphrase)
		tShouldNil(t, err)

		_, err = client.Unlockkey(a.Pubkey, _tPassphrase, nil)
		tShouldNil(t, err)
	}

	tplAddr, err := client.Addnewtemplate(AddnewtemplateParamMultisig{
		Required: 2,
		Pubkeys:  []string{a0.Pubkey, a1.Pubkey},
	})
	tShouldNil(t, err)
	tShouldNotZero(t, tplAddr)
	fmt.Println("multisig tpl addr:", *tplAddr)
	tShouldNil(t, Wait4nBlocks(3, client))

	vret, err := client.Validateaddress(*tplAddr)
	tShouldNil(t, err)
	tShouldNotZero(t, vret)
	fmt.Println("tpl addr:", toJSONIndent(*vret))

	amount := 99.0
	_, err = client.Sendfrom(CmdSendfrom{
		To:     *tplAddr,
		From:   minerAddress,
		Amount: amount,
	})
	tShouldNil(t, err)
	tShouldNil(t, Wait4balanceReach(*tplAddr, amount, client))

	fromMultisigAmount := 12.3
	_, err = client.Sendfrom(CmdSendfrom{
		To:     a2.Address,
		From:   *tplAddr,
		Amount: fromMultisigAmount,
	})
	tShouldNil(t, err)

	tShouldNil(t, Wait4balanceReach(a2.Address, fromMultisigAmount, client))

	list, err := client.Listaddress()
	tShouldNil(t, err)
	fmt.Printf("addrs: %v\n", list)
}

// 测试2个节点的多重签名(分别持有私钥1个)
func TestMultisig2Node_11(t *testing.T) {
	// 2个节点，组成网络，其中1个挖矿
	// 每个节点导入1个私钥A/B
	// 创建多签模版
	// 2个节点确保导入模版
	// 往模版转入资金
	// 创建转出交易，2个节点顺序签名
	// 广播交易，等待确认

	killCluster, nodes := TesttoolRunClusterWith2nodes(t)
	defer killCluster()

	// fmt.Println("miner:", nodes[0].MinerAddress)
	n0, n1 := nodes[0], nodes[1]
	a0, a1 := TAddr0, TAddr1

	{ //2个节点分别导入地址
		for _, imp := range []struct {
			node ClusterNode
			add  AddrKeypair
		}{{n0, a0}, {n1, a1}} {
			_, err := imp.node.Client.Importprivkey(imp.add.Privkey, _tPassphrase)
			tShouldNil(t, err)

			_, err = imp.node.Client.Unlockkey(imp.add.Pubkey, _tPassphrase, nil)
			tShouldNil(t, err)
		}
	}

	var tplAddr *string
	var err error
	{ //创建多签模版，2个节点确保导入
		tplAddr, err = n0.Client.Addnewtemplate(AddnewtemplateParamMultisig{
			Required: 2,
			Pubkeys:  []string{a0.Pubkey, a1.Pubkey},
		})
		tShouldNil(t, err)
		tShouldNotZero(t, tplAddr)
		fmt.Println("multisig tpl addr:", *tplAddr)

		vret, err := n0.Client.Validateaddress(*tplAddr)
		tShouldNil(t, err)
		tShouldNotZero(t, vret)
		fmt.Println("tpl addr:", toJSONIndent(*vret))

		importedTplAddr, err := n1.Client.Importtemplate(vret.Addressdata.Templatedata.Hex)
		tShouldNil(t, err)
		tShouldTrue(t, *tplAddr == *importedTplAddr)
	}

	amount := 99.0
	{ //往模版地址转入资金
		tShouldNil(t, Wait4balanceReach(n0.MinerAddress, 100, n0.Client))

		_, err = n0.Client.Sendfrom(CmdSendfrom{
			To:     *tplAddr,
			From:   n0.MinerAddress,
			Amount: amount,
		})
		tShouldNil(t, err)
		tShouldNil(t, Wait4balanceReach(*tplAddr, amount, n0.Client))
	}

	outFromMultisigAddr := 23.3
	{ //创建交易，分别签名，提交
		rawtx, err := n0.Client.Createtransaction(CmdCreatetransaction{
			From:   *tplAddr,
			To:     a1.Address,
			Amount: outFromMultisigAddr,
		})
		tShouldNil(t, err)
		tShouldNotZero(t, rawtx)

		sret, err := n0.Client.Signtransaction(*rawtx)
		tShouldNil(t, err)
		tShouldNotZero(t, sret)
		tShouldTrue(t, !sret.Completed)

		sret, err = n1.Client.Signtransaction(sret.Hex)
		tShouldNil(t, err)
		tShouldNotZero(t, sret)
		tShouldTrue(t, sret.Completed)

		txid, err := n1.Client.Sendtransaction(sret.Hex)
		tShouldNil(t, err)
		tShouldNotZero(t, txid)

		tShouldNil(t, Wait4balanceReach(a1.Address, outFromMultisigAddr, n1.Client))

		tx, err := n1.Client.Gettransaction(*txid, nil)
		tShouldNil(t, err)
		fmt.Println("tx from multisig", toJSONIndent(*tx))
	}

	tShouldNil(t, Wait4nBlocks(1, n1.Client))
	bal, err := n1.Client.Getbalance(nil, nil)
	tShouldNil(t, err)
	fmt.Println("节点1余额情况：", toJSONIndent(bal))
}

// 测试2个节点的多重签名(1个持有私钥，另一个只是导入模版)
func TestMultisig2Node_20(t *testing.T) {
	// 2个节点，组成网络，其中1个挖矿
	// 节点0导入私钥A/B
	// 创建多签模版
	// 2个节点确保导入模版
	// 往模版转入资金
	// 创建转出交易，节点签名
	// 广播交易，等待确认

	killCluster, nodes := TesttoolRunClusterWith2nodes(t)
	defer killCluster()

	// fmt.Println("miner:", nodes[0].MinerAddress)
	n0, n1 := nodes[0], nodes[1]
	a0, a1 := TAddr0, TAddr1

	{ //导入地址
		for _, add := range []AddrKeypair{a0, a1} {
			_, err := n0.Client.Importprivkey(add.Privkey, _tPassphrase)
			tShouldNil(t, err)
			_, err = n0.Client.Unlockkey(add.Pubkey, _tPassphrase, nil)
			tShouldNil(t, err)
		}
	}

	var tplAddr *string
	var err error
	{ //创建多签模版，2个节点确保导入
		tplAddr, err = n1.Client.Addnewtemplate(AddnewtemplateParamMultisig{
			Required: 2,
			Pubkeys:  []string{a0.Pubkey, a1.Pubkey},
		})
		tShouldNil(t, err)
		tShouldNotZero(t, tplAddr)
		fmt.Println("multisig tpl addr:", *tplAddr)

		vret, err := n1.Client.Validateaddress(*tplAddr)
		tShouldNil(t, err)
		tShouldNotZero(t, vret)
		fmt.Println("tpl addr:", toJSONIndent(*vret))

		importedTplAddr, err := n0.Client.Importtemplate(vret.Addressdata.Templatedata.Hex)
		tShouldNil(t, err)
		tShouldTrue(t, *tplAddr == *importedTplAddr)
	}

	amount := 99.0
	{ //往模版地址转入资金
		tShouldNil(t, Wait4balanceReach(n0.MinerAddress, 100, n0.Client))

		_, err = n0.Client.Sendfrom(CmdSendfrom{
			To:     *tplAddr,
			From:   n0.MinerAddress,
			Amount: amount,
		})
		tShouldNil(t, err)
		tShouldNil(t, Wait4balanceReach(*tplAddr, amount, n0.Client))
	}

	outFromMultisigAddr := 23.3
	{ //创建交易，分别签名，提交
		rawtx, err := n1.Client.Createtransaction(CmdCreatetransaction{
			From:   *tplAddr,
			To:     a1.Address,
			Amount: outFromMultisigAddr,
		})
		tShouldNil(t, err)
		tShouldNotZero(t, rawtx)

		// 钱包有问题，暂时不测decode,几个版本后看能不能用（now: 2019-12-24）
		// deTX, err := n1.Client.Decodetransaction(*rawtx)
		// tShouldNil(t, err)
		// fmt.Println("decode created tx:", toJSONIndent(deTX))

		sret, err := n0.Client.Signtransaction(*rawtx)
		tShouldNil(t, err)
		tShouldNotZero(t, sret)
		tShouldTrue(t, sret.Completed)

		txid, err := n1.Client.Sendtransaction(sret.Hex)
		tShouldNil(t, err)
		tShouldNotZero(t, txid)

		tShouldNil(t, Wait4balanceReach(a1.Address, outFromMultisigAddr, n0.Client))

		tx, err := n1.Client.Gettransaction(*txid, nil)
		tShouldNil(t, err)
		fmt.Println("tx from multisig", toJSONIndent(*tx))
	}

	tShouldNil(t, Wait4nBlocks(1, n0.Client))
	bal, err := n1.Client.Getbalance(nil, nil)
	tShouldNil(t, err)
	fmt.Println("节点1余额情况：", toJSONIndent(bal))
}

// 测试代币的发行/挖矿/交易/查询/遍历
func TestToken(t *testing.T) {
	t.Skip("当前不支持代币20191205,大概这个时间5个月后支持dpos")
}
