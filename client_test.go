// refer to: https://github.com/btcsuite/btcd/blob/master/rpcclient/infrastructure.go

package bbrpc

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	killBigBangServer, err := RunBigBangServer(&RunBigBangOptions{
		NewTmpDir: true,
		Args:      DefaultDebugBBArgs(),
	})
	tShouldNil(t, err)
	defer killBigBangServer()

	tests := []func(){
		func() { //正常获取版本
			client, err := NewClient(DefaultDebugConnConfig())
			tShouldNil(t, err)
			defer client.Shutdown()

			ver, err := client.Version()
			tShouldNil(t, err)
			tShouldTrue(t, strings.Contains(ver, "."))
			tShouldTrue(t, strings.Contains(ver, "v"))
		},
		func() { //错误的密码
			opts := DefaultDebugConnConfig()
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
