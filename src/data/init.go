package data

import (
	"errors"
	"fmt"
	"offstorage/ipfs"
	"offstorage/json_op"
	"path/filepath"

	"github.com/google/uuid"
)

type Data_api struct {
	// data存储部分所需变量
	role            string // 来自计算者or验证者的任务
	data_key_name   string // 用于访问底层数据库的密钥名字
	data_ipns_name  string // 底层数据库的ipns id
	data_local_path string // 本地底层数据库位置

	task_id string // 读写任务id

	// IPFS Shell
	ipfs_api *ipfs.Ipfs_api

	tables *DBVisitTask
}

type ModDataApi func(api *Data_api)

func NewDataShell(mod ...ModDataApi) (api *Data_api, err error) {
	api = &Data_api{
		ipfs_api: new(ipfs.Ipfs_api),
		tables:   new(DBVisitTask),
	}

	for _, fn := range mod {
		fn(api)
	}

	err = api.InitData()
	if err != nil {
		return nil, err
	}

	return api, nil
}

func DataWithRole(role string) ModDataApi {
	return func(api *Data_api) {
		api.role = role
	}
}

func DataWithKeyName(name string) ModDataApi {
	return func(api *Data_api) {
		api.data_key_name = name
	}
}

func DataWithIpnsName(name string) ModDataApi {
	return func(api *Data_api) {
		api.data_ipns_name = name
	}
}

func DataWithLocalPath(path string) ModDataApi {
	return func(api *Data_api) {
		api.data_local_path = path
	}
}

func DataWithHost(host string) ModDataApi {
	return func(api *Data_api) {
		api.ipfs_api.Ipfs_host = host
	}
}

func DataWithPort(port int) ModDataApi {
	return func(api *Data_api) {
		api.ipfs_api.Ipfs_port = port
	}
}

func (v *Data_api) InitData() (err error) {
	if v.role != "executer" && v.role != "verifier" {
		return errors.New("must give a valid role 'executer' or 'verifier'")
	}

	if v.role == "executer" {
		v.task_id = uuid.New().String()
	} else {
		if len(v.task_id) != 16 {
			return errors.New("must give a valid task id")
		}
	}

	err = v.ipfs_api.InitSh()
	if err != nil {
		return err
	}

	// 下载整个DB目录
	dest_dir := v.data_local_path
	table_ipns_path := filepath.Join("/ipns/", v.data_ipns_name)
	fmt.Printf("Download DB %s to %s \n", table_ipns_path, dest_dir)

	err = v.ipfs_api.GetFile(table_ipns_path, dest_dir)
	if err != nil {
		return err
	}

	table_dir := filepath.Join(dest_dir, "tables", v.task_id)
	if v.role == "verifier" {
		err = json_op.JsonToTable(table_dir, v.tables)
		if err != nil {
			return err
		}
	} else {
		err = json_op.GenEmptyTable(table_dir)
		if err != nil {
			return err
		}
	}

	return nil
}
