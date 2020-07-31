package bbrpc

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"testing"
)

type DockerCore struct {
	Client         *Client
	MinerAddress   string //矿工地址（资金来源)
	MinerOwnerPubk string //矿工owner address 公钥
	UnlockPass     string //已有地址解锁密码
	Conf           ConnConfig
}

// MustDockerRunDevCore 运行1个bbc core, 运行失败时t.Fatal, 自动注册 cleanup func
func MustDockerRunDevCore(t *testing.T, imageName string) DockerCore {
	killFunc, info, err := CmdRunDockerDevCore(imageName)
	if killFunc != nil {
		t.Cleanup(killFunc)
	}
	if err != nil {
		t.Fatal("unable to run bbc core dev container", err)
	}
	return info
}

// CmdRunDockerDevCore 运行1个bbc core
func CmdRunDockerDevCore(imageName string) (func(), DockerCore, error) {
	info := DockerCore{
		MinerAddress:   "20g003rgxdn4s64r4d0dchvb87p791q4epswkn1txadgv1evjqqwv70e5",
		MinerOwnerPubk: "3bc3e5f2e5e44f1cdbc44d3bf9325c93314be123f7563b8e6a88dc6eb1a25465",
		UnlockPass:     "123",
		Conf: ConnConfig{
			User:       "bbc",
			Pass:       "123",
			DisableTLS: true,
		},
	}

	idlePort, err := GetIdlePort()
	if err != nil {
		return func() {}, info, err
	}

	shellCmd := fmt.Sprintf("docker run --rm -p 9550:%d %s", idlePort, imageName)

	cmd := exec.Command(shellCmd)
	err = cmd.Run()
	if err != nil {
		return func() {}, info, err
	}
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return func() {}, info, err
	}
	log.Println("[info] docker dev core started, container id: ", string(outputBytes))

	stopContainer := func() {
		if er := exec.Command(fmt.Sprintf("docker kill %s", string(outputBytes))).Run(); er != nil {
			log.Println("[warn] stop bbc core dev container err", er)
		}
	}

	info.Conf.Host = fmt.Sprintf("127.0.0.1:%d", idlePort)
	info.Client, err = NewClient(&info.Conf)
	if err != nil {
		return stopContainer, info, err
	}
	return stopContainer, info, nil
}

// 说明：该函数代码有效，避免引入过多docker依赖注释掉了，需要的话自行拷贝使用
// DockerRunDevCore 运行1个bbc core
// func DockerRunDevCore(imageName string) (func(), DockerCore, error) {
// 	info := DockerCore{
// 		MinerAddress:   "20g003rgxdn4s64r4d0dchvb87p791q4epswkn1txadgv1evjqqwv70e5",
// 		MinerOwnerPubk: "3bc3e5f2e5e44f1cdbc44d3bf9325c93314be123f7563b8e6a88dc6eb1a25465",
// 		UnlockPass:     "123",
// 		Conf: ConnConfig{
// 			User:       "bbc",
// 			Pass:       "123",
// 			DisableTLS: true,
// 		},
// 	}
// 	cli, err := client.NewEnvClient()
// 	if err != nil {
// 		return func() {}, info, err
// 	}
// 	idlePort, err := GetIdlePort()
// 	if err != nil {
// 		return func() {}, info, err
// 	}

// 	cont, err := cli.ContainerCreate(context.Background(), &container.Config{
// 		// AttachStderr: true,
// 		// AttachStdout: true,
// 		// Tty:          true,
// 		Image:        imageName,
// 		ExposedPorts: nat.PortSet{"9550": struct{}{}},
// 	}, &container.HostConfig{
// 		PortBindings:    nat.PortMap{"9550": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: strconv.Itoa(idlePort)}}},
// 		PublishAllPorts: true,
// 		AutoRemove:      true,
// 	}, &network.NetworkingConfig{}, "")
// 	if err != nil {
// 		return func() {}, info, err
// 	}
// 	err = cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
// 	if err != nil {
// 		return func() {}, info, err
// 	}

// 	stopContainer := func() {
// 		if er := cli.ContainerStop(context.Background(), cont.ID, nil); er != nil {
// 			log.Println("[warn] stop bbc core dev container err", er)
// 		}
// 	}

// 	info.Conf.Host = fmt.Sprintf("127.0.0.1:%d", idlePort)
// 	info.Client, err = NewClient(&info.Conf)
// 	if err != nil {
// 		return stopContainer, info, err
// 	}
// 	return stopContainer, info, nil

// }

// GetIdlePort 随机获取一个空闲的端口
func GetIdlePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0") //当指定的端口为0时，操作系统会自动分配一个空闲的端口
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
