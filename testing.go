package bbrpc

import (
	"log"
	"reflect"
	"runtime/debug"
	"strings"
	"testing"
	"time"
)

func tShouldNil(t *testing.T, v interface{}, args ...interface{}) {
	if v != nil {
		debug.PrintStack()
		t.Fatalf("[test assert] should nil, but got: %v, %v", v, args)
	}
}

func tShouldTrue(t *testing.T, b bool, args ...interface{}) {
	if !b {
		debug.PrintStack()
		t.Fatalf("[test assert] should true, args: %v", args)
	}
}

func tShouldNotZero(t *testing.T, v interface{}, args ...interface{}) {
	if reflect.ValueOf(v).IsZero() {
		debug.PrintStack()
		t.Fatalf("[test assert] should not [zero value], %v", args)
	}
}

func tShouldNotContains(t *testing.T, v, containV string) {
	if strings.Contains(v, containV) {
		debug.PrintStack()
		t.Fatalf("[test assert] [%s] should not contains [%s]", v, containV)
	}
}

// ClusterNode .
type ClusterNode struct {
	IsMiner      bool
	MinerAddress string
	Client       *Client
}

//TesttoolRunClusterWith2nodes 运行2个节点，返回的第一个节点为矿工节点,发生错误则终止测试,矿工节点的日志会打印出来
func TesttoolRunClusterWith2nodes(t *testing.T) (func(), []ClusterNode) {
	killMiner, minerClient, minerAddress := TesttoolRunServerAndBeginMint(t)

	runBBOptions := DefaultDebugBBArgs()
	runBBOptions["port"] = Pstring("9901")
	runBBOptions["rpcport"] = Pstring("9907")
	ipHost := "127.0.0.1:9900"
	runBBOptions["addnode"] = &ipHost

	killPeer, err := RunBigBangServer(&RunBigBangOptions{
		NewTmpDir:       true,
		Args:            runBBOptions,
		TmpDirTag:       "peer",
		NotPrint2stdout: true,
	})
	tShouldNil(t, err)
	conn := DefaultDebugConnConfig()
	conn.Host = "127.0.0.1:9907"
	peerClient, err := NewClient(conn)
	tShouldNil(t, err, "failed to new rpc client")

	time.Sleep(time.Second)
	{ //验证2个节点确实组成了网络
		peers, err := minerClient.Listpeer()
		tShouldNil(t, err)
		tShouldTrue(t, len(peers) == 1)

		peers, err = peerClient.Listpeer()
		tShouldNil(t, err)
		tShouldTrue(t, len(peers) == 1)
	}

	return func() {
			killMiner()
			killPeer()
			minerClient.Shutdown()
			peerClient.Shutdown()
		},
		[]ClusterNode{
			{IsMiner: true, MinerAddress: minerAddress, Client: minerClient},
			{Client: peerClient},
		}
}

//TesttoolRunServerAndBeginMint （测试中）启动bigbang,使用预置的私钥开始挖矿
// 返回：killBigbang(), Client, 挖矿模版地址
func TesttoolRunServerAndBeginMint(t *testing.T) (func(), *Client, string) {
	runBBOptions := DefaultDebugBBArgs()
	runBBOptions["cryptonightaddress"] = &tCryptonightAddr.Address
	runBBOptions["cryptonightkey"] = &tCryptonightKey.Privkey

	killBigBangServer, err := RunBigBangServer(&RunBigBangOptions{
		NewTmpDir: true,
		Args:      runBBOptions,
	})
	tShouldNil(t, err, "failed to run bigbang server")

	client, err := NewClient(DefaultDebugConnConfig())
	tShouldNil(t, err, "failed to new rpc client")

	{
		_, err = client.Importprivkey(tCryptonightAddr.Privkey, _tPassphrase)
		tShouldNil(t, err)
		_, err = client.Importprivkey(tCryptonightKey.Privkey, _tPassphrase)
		tShouldNil(t, err)

		_, err = client.Unlockkey(tCryptonightAddr.Pubkey, _tPassphrase, nil)
		tShouldNil(t, err)
		_, err = client.Unlockkey(tCryptonightKey.Pubkey, _tPassphrase, nil)
		tShouldNil(t, err)
	}

	//开始挖矿
	templateAddress, err := client.Addnewtemplate(AddnewtemplateParamMint{
		Mint:  tCryptonightKey.Pubkey,
		Spent: tCryptonightAddr.Address,
	})
	tShouldNil(t, err)

	return func() {
		killBigBangServer()
		client.Shutdown()
	}, client, *templateAddress
}

// 启动bigbang-server,创建一个client,调用testFn(client)
func testClientMethod(t *testing.T, testFn func(*Client)) {
	killBigBangServer, client, _ := TesttoolRunServerAndBeginMint(t)
	defer killBigBangServer()
	testFn(client)
}

// Wait4nBlocks 每次休眠1s，直到出了n个块
func Wait4nBlocks(n int64, client *Client) error {
	count, err := client.Getblockcount(nil)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		currentCount, err := client.Getblockcount(nil)
		if err != nil {
			return err
		}
		diff := *currentCount - *count
		if diff >= n {
			log.Printf("出块已达到 %d (%d)\n", n, diff)
			return nil
		}
		log.Printf("休眠1s，等待 %d 个块,当前 %d\n", n, diff)
	}
}

// Wait4balanceReach 每次休眠1s，等待地址的余额达到
func Wait4balanceReach(addr string, balance float64, client *Client) error {
	for {
		bal, err := client.Getbalance(nil, &addr)
		if err != nil {
			return err
		}

		if len(bal) > 0 && bal[0].Avail >= balance {
			log.Printf("地址 %s 余额达到%v (%v)\n", addr, balance, bal[0].Avail)
			return nil
		}

		log.Printf("休眠1s，等待地址 %s 余额达到%v (当前 %v)\n", addr, balance, bal)
		time.Sleep(time.Second)
	}
}
