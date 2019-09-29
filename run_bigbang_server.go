package bbrpc

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// debug BigBang args
func defaultDebugBBArgs() map[string]*string {
	return map[string]*string{
		"testnet":     nil,
		"listen4":     nil,
		"debug":       nil,
		"rpcport":     ps("9906"),
		"rpcuser":     ps("rpcusr"),
		"rpcpassword": ps("pwd"),
	}
}
func defaultDebugConnConfig() *ConnConfig {
	return &ConnConfig{
		Host:       "127.0.0.1:9906",
		DisableTLS: true,
		User:       "rpcusr",
		Pass:       "pwd",
	}
}

// RunBigBangOptions .
type RunBigBangOptions struct {
	NewTmpDir       bool               //创建并使用新的临时目录作为datadir
	Args            map[string]*string //k-v ,v 为nil时为flag
	NotPrint2stdout bool               //不打印到stdout
}
type killHook func() error

// RunBigBangServer run bigbang server,print out to stdout, require bigbang in the $PATH, this func is used for testing BigBang in local test env
// return func() to kill bigbang server
// usage:
// 		killBigbang, err := RunBigBangServer(options)
//  	defer killBigBang()
func RunBigBangServer(optionsPtr *RunBigBangOptions) (func(), error) {
	killHooks := []killHook{}

	var options RunBigBangOptions
	var err error

	if optionsPtr == nil {
		options = RunBigBangOptions{}
	} else {
		options = *optionsPtr
	}
	if options.Args == nil {
		options.Args = map[string]*string{}
	}

	var dataDir string
	if options.NewTmpDir {
		for k, v := range options.Args {
			if k == "datadir" {
				return nil, fmt.Errorf("datadir specified in args (%v), NewTmpDir not work", v)
			}
		}

		tmpDir := strings.TrimRight(os.TempDir(), "/")
		dataDir = tmpDir + "/bigbang_data_tmp_" + time.Now().Format(time.RFC3339) + "/"
		err := os.MkdirAll(dataDir, 0777)
		if err != nil {
			return nil, fmt.Errorf("cannot create tmp dir: %v, err: %v", dataDir, err)
		}
		options.Args["datadir"] = &dataDir

		killHooks = append(killHooks, func() error {
			return os.RemoveAll(dataDir)
		})
	}

	args := []string{}
	for k, v := range options.Args {
		arg := "-" + k
		if v != nil {
			arg += "=" + *v
		}
		args = append(args, arg)
	}

	closeChan := make(chan struct{})

	cmd := exec.Command("bigbang", args...)
	fmt.Println("[debug] bigbang-server args", cmd.Args)
	if options.NotPrint2stdout {
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	go func() {
		fmt.Println("Waiting for message to kill bigbang")
		<-closeChan
		fmt.Println("Received message,killing bigbang server")

		if e := cmd.Process.Kill(); e != nil {
			fmt.Println("关闭 bigbang 时发生异常", e)
		}
		closeChan <- struct{}{}
	}()

	// err = cmd.Wait()
	fmt.Println("等待1秒,让 bigbang 启动")
	time.Sleep(time.Millisecond * 1000)
	return func() {
		closeChan <- struct{}{}
		<-closeChan
	}, nil
}
