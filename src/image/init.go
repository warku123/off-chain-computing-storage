package image

import (
	"fmt"
	"offstorage/ipfs"
	"offstorage/json_op"
	"path/filepath"

	"github.com/cbergoon/merkletree"
)

type Image_api struct {
	ipfs_api *ipfs.Ipfs_api

	// image存储部分所需变量
	task_name        string // 任务名，用于索引对应快照
	image_key_name   string
	image_ipns_name  string
	image_local_path string

	// Hash状态索引表
	image_table *ImageTable
	// 待优化，有点冗余，按理说只要拿出来对应task_name索引的[]Image就可以了

	// Merkel树
	merkle_tree *merkletree.MerkleTree
}

type ModImageApi func(api *Image_api)

func NewImageShell(mod ...ModImageApi) (api *Image_api, err error) {
	api = &Image_api{
		ipfs_api:    new(ipfs.Ipfs_api),
		image_table: new(ImageTable),
		merkle_tree: new(merkletree.MerkleTree),
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

func ImageWithTaskName(name string) ModImageApi {
	return func(api *Image_api) {
		api.task_name = name
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

	// download the image table from ipns
	dest_dir := v.image_local_path
	table_ipns_path := filepath.Join("/ipns/", v.image_ipns_name)
	fmt.Printf("Download ImageTable %s to %s \n", table_ipns_path, dest_dir)
	err = v.ipfs_api.GetFile(table_ipns_path, dest_dir)
	if err != nil {
		return err
	}

	local_imagetable_path := filepath.Join(dest_dir, v.image_ipns_name)
	err = json_op.JsonToTable(local_imagetable_path, v.image_table)
	if err != nil {
		return err
	}

	// build the merkle tree
	err = v.BuildTree()
	if err != nil {
		return err
	}

	return nil
}

func (v *Image_api) CloseImage() error {
	return nil
}
