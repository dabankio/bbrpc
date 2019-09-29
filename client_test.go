// refer to: https://github.com/btcsuite/btcd/blob/master/rpcclient/infrastructure.go

package bbrpc

import (
	"fmt"
	"strings"
	"testing"
)

// 启动bigbang-server,创建一个client,调用testFn(client)
func testClientMethod(t *testing.T, testFn func(*Client)) {
	killBigBangServer, err := RunBigBangServer(&RunBigBangOptions{
		NewTmpDir: true,
		Args:      defaultDebugBBArgs(),
	})
	tShouldNil(t, err, "failed to run bigbang server")
	defer killBigBangServer()

	client, err := NewClient(defaultDebugConnConfig())
	tShouldNil(t, err, "failed to new rpc client")
	defer client.Shutdown()

	testFn(client)
}

func TestNewClient(t *testing.T) {
	killBigBangServer, err := RunBigBangServer(&RunBigBangOptions{
		NewTmpDir: true,
		Args:      defaultDebugBBArgs(),
	})
	tShouldNil(t, err)
	defer killBigBangServer()

	tests := []func(){
		func() { //正常获取版本
			client, err := NewClient(defaultDebugConnConfig())
			tShouldNil(t, err)
			defer client.Shutdown()

			ver, err := client.Version()
			tShouldNil(t, err)
			tShouldTrue(t, strings.Contains(ver, "."))
			tShouldTrue(t, strings.Contains(ver, "v"))
		},
		func() { //错误的密码
			opts := defaultDebugConnConfig()
			opts.Pass = "bad_pass"

			c, err := NewClient(opts)
			tShouldNil(t, err)
			defer c.Shutdown()

			_, err = c.Version()
			tShouldTrue(t, err != nil)
			tShouldTrue(t, strings.Contains(err.Error(), "401"), "not contains 401")
			fmt.Println("bad pass error:", err.Error())
		},
	}

	for _, fn := range tests {
		fn()
	}
}
