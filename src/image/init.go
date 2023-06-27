package image

import (
	"offstorage/ipfs"
)

type Image_api struct {
	ipfs_api *ipfs.Ipfs_api

	// image存储部分所需变量
	chain_name       string // 链名，用于标示快照名字&本地存储快照索引的文件夹父目录名
	image_key_name   string
	image_ipns_name  string
	image_local_path string
}

type ModImageApi func(api *Image_api)

func NewImageShell(mod ...ModImageApi) (api *Image_api, err error) {
	api = &Image_api{}
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

func (v *Image_api) initImage()
