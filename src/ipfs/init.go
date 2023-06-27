package ipfs

import (
	"fmt"

	shell "github.com/ipfs/go-ipfs-api"
)

type Ipfs_api struct {
	// 初始化ipfs部分
	Ipfs_host string
	Ipfs_port int

	Sh *shell.Shell
}

type ModIpfsApi func(api *Ipfs_api)

func NewShell(mod ...ModIpfsApi) (*Ipfs_api, error) {
	api := &Ipfs_api{
		Ipfs_host: "127.0.0.1",
		Ipfs_port: 5001,
	}

	for _, fn := range mod {
		fn(api)
	}

	err := api.InitSh()
	if err != nil {
		return nil, err
	}

	fmt.Println("New IPFS Shell created!")
	return api, nil
}

func (v *Ipfs_api) InitSh() error {
	// 创建IPFS访问
	v.Sh = shell.NewShell(fmt.Sprintf("%s:%d", v.Ipfs_host, v.Ipfs_port))

	return nil
}

func ShellWithHost(host string) ModIpfsApi {
	return func(api *Ipfs_api) {
		api.Ipfs_host = host
	}
}

func ShellWithPort(port int) ModIpfsApi {
	return func(api *Ipfs_api) {
		api.Ipfs_port = port
	}
}
