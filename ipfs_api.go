package storage_off

import (
	"errors"
	"fmt"
	"os"

	shell "github.com/ipfs/go-ipfs-api"
)

type ipfs_api struct {
	ipfs_host string
	ipfs_port int

	snapshot_tag string

	sh *shell.Shell
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

	err := api.initSh()
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

func (v *ipfs_api) initSh() error {
	v.sh = shell.NewShell(fmt.Sprintf("%s:%d", v.ipfs_host, v.ipfs_port))

	if len(v.snapshot_tag) == 0 {
		return errors.New("Must have a dir.")
	}

	if v.snapshot_tag[0] != '/' {
		return errors.New("Dir must begin with \"/\".")
	}

	dest_dir := v.snapshot_tag
	dir_stat, err := os.Stat(dest_dir)
	if err != nil {
		return err
	}

	dir_exist := dir_stat.IsDir()

	if !dir_exist {
		err := os.MkdirAll(v.snapshot_tag, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
