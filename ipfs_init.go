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

	snapshot_tag string

	sh      *shell.Shell
	key     *shell.Key
	ipns_id string
}

type ModIpfsApi func(api *ipfs_api)

// NewShell to create a new IPFS interface
// usage: NewShell(ShellWithHost, ShellWithPort...)
func NewShell(mod ...ModIpfsApi) (*ipfs_api, error) {
	api := ipfs_api{
		ipfs_host:    "127.0.0.1",
		ipfs_port:    5001,
		snapshot_tag: "/ipfs_files",
	}

	for _, fn := range mod {
		fn(&api)
	}

	_, err := api.initSh()
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

// return IPNS id
func (v *ipfs_api) initSh() (string, error) {
	v.sh = shell.NewShell(fmt.Sprintf("%s:%d", v.ipfs_host, v.ipfs_port))

	if len(v.snapshot_tag) == 0 {
		return "", errors.New("must have a dir\n")
	}

	if v.snapshot_tag[0] != '/' {
		return "", errors.New("dir must begin with \"/\"\n")
	}

	dest_dir := v.snapshot_tag
	dir_stat, err := os.Stat(dest_dir)
	if err != nil {
		return "", err
	}
	dir_exist := dir_stat.IsDir()
	if !dir_exist {
		err := os.MkdirAll(v.snapshot_tag, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	table_dest_dir := fmt.Sprintf("%s/snapshot.txt", v.snapshot_tag)
	_, err = os.Stat(table_dest_dir)
	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(table_dest_dir)
			if err != nil {
				return "", err
			}
		}
	}

	err = v.genIPNSkey()
	if err != nil {
		return "", err
	}

	cid, err := v.AddFolder(table_dest_dir)
	if err != nil {
		return "", err
	}

	v.ipns_id, err = v.PublishFile(cid)
	if err != nil {
		return "", err
	}

	return v.ipns_id, err
}

func (v *ipfs_api) genIPNSkey() (err error) {
	v.key, err = v.sh.KeyGen(context.Background(), v.snapshot_tag)
	if err != nil {
		return err
	}
	return nil
}
