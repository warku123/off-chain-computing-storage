package data

import (
	"errors"
	"fmt"
	"offstorage/json_op"
	"path/filepath"
)

func (v *Data_api) GetDataFromIPFS() (err error) {
	// 下载整个DB目录
	dest_dir := v.data_local_path
	table_ipns_path := filepath.Join("/ipns/", v.data_ipns_name)
	fmt.Printf("Download DB %s to %s \n", table_ipns_path, dest_dir)

	err = v.ipfs_api.GetFile(table_ipns_path, dest_dir)
	if err != nil {
		return err
	}
	return nil
}

func (v *Data_api) SyncDataToIPFS() (ipns_name string, err error) {
	// 上传整个DB目录
	src_dir := v.data_local_path
	table_local_path := filepath.Join(src_dir, v.data_ipns_name)
	table_ipns_path := filepath.Join("/ipns/", v.data_ipns_name)
	fmt.Printf("Upload DB %s to %s \n", src_dir, table_ipns_path)

	// Upload db
	cid, err := v.ipfs_api.AddFile(table_local_path)
	if err != nil {
		return "", err
	}

	// Publish db
	ipns_name, err = v.ipfs_api.PublishFile(cid, v.data_key_name)
	if err != nil {
		return "", err
	}

	return ipns_name, nil
}

func (v *Data_api) GetDataCid(name string) (cid string, err error) {
	// Data in write table
	if entry, ok := v.tables.Write_table[name]; ok {
		return entry[len(entry)-1], nil
	}

	// Data in read table
	if entry, ok := v.tables.Read_table[name]; ok {
		return string(entry), nil
	}

	// Neither in read table nor in write table
	if v.role == "verifier" {
		return "", errors.New("get data: no data in verfier's read table")
	}

	// Data in db
	err = v.GetDataFromIPFS()
	if err != nil {
		return "", err
	}

	db_dir := filepath.Join(v.data_local_path, v.data_ipns_name, "db")
	err = json_op.JsonToTable(db_dir, v.db)
	if err != nil {
		return "", err
	}

	data_length, err := v.db.GetDataVersionNum(name)
	if err != nil {
		return "", err
	}
	cid = v.db.GetCid(name, data_length-1)
	v.db.AddReadNum(name, data_length-1)

	v.tables.AddReadTuple(name, cid)

	// 可以并行，待优化
	_, err = v.SyncDataToIPFS()
	if err != nil {
		return "", err
	}

	return cid, nil
}

func (v *Data_api) GetData(name string) (value string, err error) {
	cid, err := v.GetDataCid(name)
	if err != nil {
		return "", err
	}
	value, err = v.ipfs_api.ReadFile("/ipfs/" + cid)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (v *Data_api) AddData(name string, value string) (err error) {
	cid, err := v.ipfs_api.AddString(value)
	if err != nil {
		return err
	}

	// Add data to write table
	v.tables.AddWriteTuple(name, cid)

	return nil
}
