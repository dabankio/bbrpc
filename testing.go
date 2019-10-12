package bbrpc

import (
	"reflect"
	"runtime/debug"
	"strings"
	"testing"
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

// （测试中）启动bigbang,使用预置的私钥开始挖矿
// 返回：killBigbang(), Client, 挖矿模版地址
func tRunBigbangServerAndBeginMint(t *testing.T) (func(), *Client, string) {
	runBBOptions := DefaultDebugBBArgs()
	runBBOptions["cryptonightaddress"] = &tCryptonightAddr.Address
	runBBOptions["cryptonightkey"] = &tCryptonightKey.Privkey

	killBigBangServer, err := RunBigBangServer(&RunBigBangOptions{
		NewTmpDir: true,
		Args:      runBBOptions,
	})
	tShouldNil(t, err, "failed to run bigbang server")
	// defer killBigBangServer()

	client, err := NewClient(DefaultDebugConnConfig())
	tShouldNil(t, err, "failed to new rpc client")
	// defer client.Shutdown()

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
