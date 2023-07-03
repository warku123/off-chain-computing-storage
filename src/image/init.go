package image

import (
	"fmt"
	"offstorage/ipfs"
)

type Image_api struct {
	ipfs_api *ipfs.Ipfs_api

	// image存储部分所需变量
	chain_name       string // 链名，用于标示快照名字&本地存储快照索引的文件夹父目录名
	image_key_name   string
	image_ipns_name  string
	image_local_path string

	// Hash状态索引表
	image_table *ImageTable
}

type ModImageApi func(api *Image_api)

func NewImageShell(mod ...ModImageApi) (api *Image_api, err error) {
	api = &Image_api{
		ipfs_api: new(ipfs.Ipfs_api),
	}

	for _, fn := range mod {
		fn(api)
	}

	err = api.InitImage()
	if err != nil {
		return nil, err
	}

	return api, nil
}

func ImageWithChainName(name string) ModImageApi {
	return func(api *Image_api) {
		api.chain_name = name
	}
}

func ImageWithKeyName(name string) ModImageApi {
	return func(api *Image_api) {
		api.image_key_name = name
	}
}

func ImageWithIpnsName(name string) ModImageApi {
	return func(api *Image_api) {
		api.image_ipns_name = name
	}
}

func ImageWithLocalPath(path string) ModImageApi {
	return func(api *Image_api) {
		api.image_local_path = path
	}
}

func ImageWithHost(host string) ModImageApi {
	return func(api *Image_api) {
		api.ipfs_api.Ipfs_host = host
	}
}

func ImageWithPort(port int) ModImageApi {
	return func(api *Image_api) {
		api.ipfs_api.Ipfs_port = port
	}
}

func (v *Image_api) InitImage() error {
	defer fmt.Println("New Image Shell created!")
	err := v.ipfs_api.InitSh()
	if err != nil {
		return err
	}

	dest_dir := v.image_local_path
	table_ipns_path := fmt.Sprintf("/ipns/%s", v.image_ipns_name)

	fmt.Printf("Download ImageTable %s to %s /n", table_ipns_path, dest_dir)

	err = v.ipfs_api.GetFile(table_ipns_path, dest_dir)
	if err != nil {
		return err
	}

	return nil
}
