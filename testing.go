package bbrpc

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strconv"
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
	RPCPort      int
}

//TesttoolRunClusterWith2nodes 运行2个节点，返回的第一个节点为矿工节点,发生错误则终止测试,矿工节点的日志会打印出来
func TesttoolRunClusterWith2nodes(t *testing.T) (func(), []ClusterNode) {
	killMiner, minerClient, minerAddress := TesttoolRunServerAndBeginMint(t)

	runBBOptions := DefaultDebugBBArgs()

	runBBOptions["port"] = Pstring(strconv.Itoa(TDefaultPort2))
	runBBOptions["rpcport"] = Pstring(strconv.Itoa(TDefaultRPCPort2))
	ipHost := "127.0.0.1:" + strconv.Itoa(TDefaultPort)
	runBBOptions["addnode"] = &ipHost

	killPeer, err := RunBigBangServer(&RunBigBangOptions{
		NewTmpDir:       true,
		Args:            runBBOptions,
		TmpDirTag:       "peer",
		NotPrint2stdout: true,
	})
	tShouldNil(t, err)
	conn := DefaultDebugConnConfig()
	conn.Host = "127.0.0.1:" + strconv.Itoa(TDefaultRPCPort2)
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
			{IsMiner: true, MinerAddress: minerAddress, Client: minerClient, RPCPort: TDefaultRPCPort},
			{Client: peerClient, RPCPort: TDefaultRPCPort2},
		}
}

//TesttoolRunServerAndBeginMint （测试中）启动bigbang,使用预置的私钥开始挖矿, opts 作为可选参数只使用第一个值(如果有)
// 返回：killBigbang(), Client, 挖矿模版地址
func TesttoolRunServerAndBeginMint(t *testing.T, opts ...RunBigBangOptions) (func(), *Client, string) {
	runBBOptions := DefaultDebugBBArgs()
	runBBOptions["cryptonightaddress"] = &TCryptonightAddr.Address
	runBBOptions["cryptonightkey"] = &TCryptonightKey.Privkey

	opt := RunBigBangOptions{
		NewTmpDir: true,
		Args:      runBBOptions,
	}
	if len(opts) > 0 {
		opt = opts[0]
		if len(opt.Args) == 0 {
			opt.Args = runBBOptions
		} else {
			for k, v := range runBBOptions { //补充没有的参数
				if _, ok := opt.Args[k]; !ok {
					opt.Args[k] = v
				}
			}
		}
	}
	if !testing.Verbose() {
		opt.NotPrint2stdout = true
	}
	killBigBangServer, err := RunBigBangServer(&opt)
	tShouldNil(t, err, "failed to run bigbang server")

	client, err := NewClient(DefaultDebugConnConfig())
	tShouldNil(t, err, "failed to new rpc client")

	{
		_, _ = client.Importprivkey(TCryptonightAddr.Privkey, _tPassphrase, nil)
		_, _ = client.Importprivkey(TCryptonightKey.Privkey, _tPassphrase, nil) //这个无需导入，配置已有，导入反而报错
		_, _ = client.Unlockkey(TCryptonightAddr.Pubkey, _tPassphrase, nil)
		_, _ = client.Unlockkey(TCryptonightKey.Pubkey, _tPassphrase, nil)
	}

	//开始挖矿
	templateAddress, err := client.Addnewtemplate(AddnewtemplateParamMint{
		Mint:  TCryptonightKey.Pubkey,
		Spent: TCryptonightAddr.Address,
	})
	tShouldNil(t, err)

	return func() {
		killBigBangServer()
		client.Shutdown()
	}, client, *templateAddress
}

// 启动bigbang-server,创建一个client,调用testFn(client)
func testClientMethod(t *testing.T, testFn func(*Client)) {
	// killBigBangServer, client, _ := TesttoolRunServerAndBeginMint(t)
	// defer killBigBangServer()
	// testFn(client)

	info := MustDockerRunDevCore(t, "bbccore:0.4")
	testFn(info.Client)
}

// 启动bigbang-server,创建一个client,调用testFn(client)
func runClientTest(t *testing.T, testFn func(*Client, string)) {
	killBigBangServer, client, minerAddr := TesttoolRunServerAndBeginMint(t)
	defer killBigBangServer()
	testFn(client, minerAddr)
}

// Wait4nBlocks 每次休眠1s，直到出了n个块
func Wait4nBlocks(n int, client *Client) error {
	count, err := client.Getblockcount(nil)
	if err != nil {
		return err
	}

	fmt.Printf("等待 %d 个块 ", n)
	prevDiff := 0
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
			fmt.Printf("[达到(%d)]\n", diff)
			return nil
		}
		fmt.Print(".")
		if prevDiff != diff {
			fmt.Print(diff)
			prevDiff = diff
		}
	}
}

// Wait4balanceReach 每次休眠1s，等待地址的余额达到
func Wait4balanceReach(addr string, balance float64, client *Client, msg ...string) error {
	fmt.Printf("%v等待地址 %s 余额达到%v ", msg, addr, balance)

	prevBal := 0.0
	for {
		bal, err := client.Getbalance(nil, &addr)
		if err != nil {
			return err
		}

		f := 0.0
		if len(bal) > 0 {
			f = bal[0].Avail - bal[0].Unconfirmed
		}

		fmt.Printf(".")
		if f != prevBal {
			prevBal = f
			fmt.Printf("%v", f)
		}

		if f >= balance {
			fmt.Printf("[达到]\n")
			return nil
		}
		time.Sleep(time.Second)
	}
}
