package image

import (
	"offstorage/ipfs"
)

type Image_api struct {
	ipfs_api *ipfs.Ipfs_api

	// image存储部分所需变量
	chain_name       string // 链名，用于标示快照名字&本地存储快照索引的文件夹父目录名
	image_key_name   string
	image_ipns_id    string
	image_local_path string
}

type ModImageApi func(api *Image_api)

func NewImageShell(mod ...ModImageApi) (api *Image_api, error) {
	api = &Image_api{}
	
}

func (v *Image_api) initImage()
