package image

import (
	"context"
	"errors"
	"fmt"

	shell "github.com/ipfs/go-ipfs-api"
)

type ipfs_api struct {
	// 初始化ipfs部分
	ipfs_host string
	ipfs_port int

	sh *shell.Shell

	// image存储部分所需变量
	snapshot_tag     string // 链名，用于标示快照名字&本地存储快照索引的文件夹父目录名
	image_key_name   string
	image_ipns_id    string
	image_local_path string

	image_idx int // image的索引，自动生成，用于标识DB中的读写变量表归属于哪个任务

	// data存储部分所需变量
	role            string // 来自计算者or验证者的任务
	data_key_name   string // 用于访问底层数据库的密钥名字
	data_ipns_id    string // 底层数据库的ipns id
	data_local_path string // 本地底层数据库位置

	data_visit_task *DBVisitTask
}

type ModIpfsApi func(api *ipfs_api)

// NewShell to create a new IPFS interface
// usage: NewShell(ShellWithHost, ShellWithPort...)
func NewShell(mod ...ModIpfsApi) (*ipfs_api, error) {
	api := ipfs_api{
		ipfs_host:        "127.0.0.1",
		ipfs_port:        5001,
		image_local_path: "/images",
		data_local_path:  "/off_chain_data",
	}

	for _, fn := range mod {
		fn(&api)
	}

	_, _, err := api.initSh()
	if err != nil {
		return nil, err
	}

	fmt.Println("New IPFS Shell created!")
	return &api, nil
}

func ShellWithHost(host string) ModIpfsApi {
	return func(api *ipfs_api) {
		api.ipfs_host = host
	}
}

func ShellWithPort(port int) ModIpfsApi {
	return func(api *ipfs_api) {
		api.ipfs_port = port
	}
}

func ShellWithRole(role string) ModIpfsApi {
	return func(api *ipfs_api) {
		api.role = role
	}
}

// Chain name in this visit
func ShellWithTag(tag string) ModIpfsApi {
	return func(api *ipfs_api) {
		api.snapshot_tag = tag
	}
}

func ShellWithDataKeyName(key string) ModIpfsApi {
	return func(api *ipfs_api) {
		api.data_key_name = key
	}
}

func ShellWithImageKeyName(key string) ModIpfsApi {
	return func(api *ipfs_api) {
		api.image_key_name = key
	}
}

// return IPNS id
func (v *ipfs_api) initSh() (string, string, error) {
	if len(v.snapshot_tag) == 0 {
		return "", "", errors.New("must have a snapshot tag\n")
	}

	if len(v.image_key_name) == 0 {
		return "", "", errors.New("must have a key to visit image DB\n")
	}

	if len(v.data_key_name) == 0 {
		return "", "", errors.New("must have a key to visit data DB\n")
	}

	if v.role != "executer" && v.role != "verifier" {
		return "", "", errors.New("must give a valid role 'executer' or 'verifier'\n")
	}

	// 创建IPFS访问
	v.sh = shell.NewShell(fmt.Sprintf("%s:%d", v.ipfs_host, v.ipfs_port))

	err := v.InitImage()
	if err != nil {
		return "", "", err
	}

	// err = v.InitDBVisit()
	// if err != nil {
	// 	return "", "", err
	// }

	return v.image_key_name, v.image_ipns_id, err
}

func (v *ipfs_api) CloseSh(image string) (err error) {
	err = v.EndDBVisit()
	if err != nil {
		return err
	}

	_, _, err = v.NewImage(image)
	if err != nil {
		return err
	}

	return nil
}

func (v *ipfs_api) genIPNSkey() (key *shell.Key, err error) {
	key, err = v.sh.KeyGen(context.Background(), v.snapshot_tag)
	if err != nil {
		return nil, err
	}
	return key, nil
}
