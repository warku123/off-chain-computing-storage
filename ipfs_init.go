package storage_off

import (
	"context"
	"errors"
	"fmt"
	"os"

	shell "github.com/ipfs/go-ipfs-api"
)

type ipfs_api struct {
	ipfs_host string
	ipfs_port int

	snapshot_tag  string // 标示快照名字&ipns中用来发布该快照的密钥名字&本地存储快照索引的文件夹父目录名
	image_key     *shell.Key
	image_ipns_id string

	role            string // 来自计算者or验证者的任务
	data_key_name   string // 用于访问底层数据库的密钥名字
	data_ipns_id    string // 底层数据库的ipns id
	data_local_path string // 本地底层数据库位置

	sh *shell.Shell

	data_visit_task *DBVisitTask
}

type ModIpfsApi func(api *ipfs_api)

// NewShell to create a new IPFS interface
// usage: NewShell(ShellWithHost, ShellWithPort...)
func NewShell(mod ...ModIpfsApi) (*ipfs_api, error) {
	api := ipfs_api{
		ipfs_host:       "127.0.0.1",
		ipfs_port:       5001,
		snapshot_tag:    "/default_chain",
		data_local_path: "/off_chain_data",
	}

	for _, fn := range mod {
		fn(&api)
	}

	_, _, err := api.initSh()
	if err != nil {
		return nil, err
	}
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

// IPFS Files system file name
func ShellWithDirTag(tag string) ModIpfsApi {
	return func(api *ipfs_api) {
		api.snapshot_tag = tag
	}
}

func ShellWithDBKeyName(key string) ModIpfsApi {
	return func(api *ipfs_api) {
		api.data_key_name = key
	}
}

// return IPNS id
func (v *ipfs_api) initSh() (string, string, error) {
	if len(v.snapshot_tag) == 0 {
		return "", "", errors.New("must have a dir\n")
	}

	if v.snapshot_tag[0] != '/' {
		return "", "", errors.New("dir must begin with \"/\"\n")
	}

	if len(v.data_key_name) == 0 {
		return "", "", errors.New("must have a key to visit DB\n")
	}

	if v.role != "executer" || v.role != "verifier" {
		return "", "", errors.New("must give a correct role\n")
	}

	v.sh = shell.NewShell(fmt.Sprintf("%s:%d", v.ipfs_host, v.ipfs_port))

	dest_dir := v.snapshot_tag
	dir_stat, err := os.Stat(dest_dir)
	if err != nil {
		return "", "", err
	}
	dir_exist := dir_stat.IsDir()
	if !dir_exist {
		err := os.MkdirAll(v.snapshot_tag, os.ModePerm)
		if err != nil {
			return "", "", err
		}
	}

	table_dest_dir := fmt.Sprintf("%s/snapshot.txt", v.snapshot_tag)
	_, err = os.Stat(table_dest_dir)
	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(table_dest_dir)
			if err != nil {
				return "", "", err
			}
		}
	}

	err = v.genIPNSkey()
	if err != nil {
		return "", "", err
	}

	// image_ipns_cid 对应的是快照索引表
	cid, err := v.AddFile(table_dest_dir)
	if err != nil {
		return "", "", err
	}

	v.image_ipns_id, err = v.PublishImage(cid)
	if err != nil {
		return "", "", err
	}
	// 这边key.Id想返回key的内容，但是不知道key里面是不是，待测试
	return v.image_key.Id, v.image_ipns_id, err
}

func (v *ipfs_api) genIPNSkey() (err error) {
	v.image_key, err = v.sh.KeyGen(context.Background(), v.snapshot_tag)
	if err != nil {
		return err
	}
	return nil
}
