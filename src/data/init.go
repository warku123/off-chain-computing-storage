package data

import (
	"errors"
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

	task_id   string // 计算者任务id
	v_task_id string // 验证者任务id

	// IPFS Shell
	ipfs_api *ipfs.Ipfs_api

	// Read and Write table
	tables *DBVisitTask

	// Whole data table
	db *Data_table
}

type ModDataApi func(api *Data_api)

func NewDataShell(mod ...ModDataApi) (api *Data_api, err error) {
	api = &Data_api{
		ipfs_api: new(ipfs.Ipfs_api),
		tables:   new(DBVisitTask),
		db:       new(Data_table),
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

func DataWithTaskID(id string) ModDataApi {
	return func(api *Data_api) {
		api.task_id = id
	}
}

func (v *Data_api) InitData() (err error) {
	if v.role != "executer" && v.role != "verifier" && v.role != "judge" {
		return errors.New("must give a valid role 'executer'/'verifier'/'judge'")
	}

	if v.role == "executer" {
		v.task_id = uuid.New().String()
	} else {
		if len(v.task_id) != 16 {
			return errors.New("must give a valid task id")
		}
		if v.role == "verifier" {
			v.v_task_id = uuid.New().String()
		}
	}

	// 初始化ipfs shell
	err = v.ipfs_api.InitSh()
	if err != nil {
		return err
	}

	// 下载整个Data存储
	err = v.GetDataFromIPFS()
	if err != nil {
		return err
	}

	// 读取tables
	table_dir := filepath.Join(v.data_local_path, v.data_ipns_name, "executer", v.task_id)
	if v.role == "executer" {
		err = json_op.GenEmptyTable(table_dir)
		if err != nil {
			return err
		}
	} else {
		// 读读写表
		err = json_op.JsonToTable(table_dir, v.tables)
		if err != nil {
			return err
		}
		if v.role == "verifier" {
			// 清空写表
			v.tables.Write_table = make(map[string]write_variable)
		}
	}

	// 读取db
	db_dir := filepath.Join(v.data_local_path, v.data_ipns_name, "db")
	err = json_op.JsonToTable(db_dir, v.db)
	if err != nil {
		return err
	}

	return nil
}

func (v *Data_api) CloseImage() (err error) {
	// 链上judge的逻辑还没加
	var table_path string
	if v.role == "executer" {
		table_path = filepath.Join(v.data_local_path, v.data_ipns_name, "executer", v.task_id)
	} else {
		table_path = filepath.Join(v.data_local_path, v.data_ipns_name, "verifier", v.task_id+"_"+v.v_task_id)
	}
	json_bytes, err := json_op.TableToJson(v.tables)
	if err != nil {
		return err
	}

	err = json_op.SaveJsonToFile(table_path, json_bytes)
	if err != nil {
		return err
	}

	// 上传table
	_, err = v.SyncDataToIPFS()
	if err != nil {
		return err
	}

	return nil
}
