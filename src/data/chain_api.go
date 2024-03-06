package data

import (
	"errors"
	"offstorage/utils"
	"os"
	"path/filepath"
)

func (v *Data_api) DeleteTables() (err error) {
	// 这块应该也要加锁
	if v.role != "judge" {
		return errors.New("only judge can delete tables")
	}

	exe_table_dir := filepath.Join(v.data_local_path, "executer", v.task_id)
	err = os.Remove(exe_table_dir)
	if err != nil {
		return err
	}

	ver_table_dir := filepath.Join(v.data_local_path, "verifier", v.task_id)
	if utils.PathExists(ver_table_dir) {
		err = os.RemoveAll(ver_table_dir)
		if err != nil {
			return err
		}
	}

	// 删完立即同步IPFS
	_, err = v.SyncDataToIPFS()
	if err != nil {
		return err
	}

	return nil
}

// DataPersistance 用于根据最初的task_id，持久化数据到DB
func (v *Data_api) DataPersistance() (err error) {
	if v.role != "judge" {
		return errors.New("only judge can persist data")
	}

	for key, entry := range v.tables.Write_table {
		cid := entry[len(entry)-1]
		// fmt.Println(1111)
		err = v.db.AddCid(key, cid)
		if err != nil {
			return err
		}
		err = v.db.ReduceWriteNum(key)
		if err != nil {
			return err
		}
	}

	for key, entry := range v.tables.Read_table {
		err = v.db.ReduceReadNum(key, entry.Read_version)
		if err != nil {
			return err
		}
	}

	// 删除表应该和持久化数据一起，否则会产生数据不一致
	err = v.DeleteTables()
	if err != nil {
		return err
	}

	// 修改完读写表顺带gc掉没有读任务的数据
	err = v.DBGarbageCollection()
	if err != nil {
		return err
	}

	return nil
}

func (v *Data_api) ReduceReadNum() (err error) {
	if v.role != "judge" {
		return errors.New("only judge can reduce read num")
	}

	for key, entry := range v.tables.Read_table {
		err = v.db.ReduceReadNum(key, entry.Read_version)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *Data_api) DBGarbageCollection() (err error) {
	if v.role != "judge" {
		return errors.New("only judge can do garbage collection")
	}

	for key, entry := range *(v.db) {
		var i int
		// 检查到第一个read_num不为0的数据，否则会造成读取错误
		last := entry.Data_tuples[len(entry.Data_tuples)-1]
		for i = 0; i < len(entry.Data_tuples)-1; i++ {
			if entry.Data_tuples[i].Read_num == 0 {
				entry.Gc_offset++
			} else {
				break
			}
		}
		new_data_list := entry.Data_tuples[:i]
		entry.Data_tuples = new_data_list
		entry.Data_tuples = append(entry.Data_tuples, last)
		(*(v.db))[key] = entry
	}
	return nil
}

func (v *Data_api) InitDB() (err error) {
	if v.role != "judge" {
		return errors.New("only judge can init db")
	}

	// v.db = new(Data_table)
	err = utils.GenEmptyTable(
		filepath.Join(v.data_local_path, v.data_ipns_name, "db"),
	)
	if err != nil {
		return err
	}

	err = v.DataPersistance()
	if err != nil {
		return err
	}
	return nil
}

func (v *Data_api) TraverseTable(task_id, v_task_id string) (table string, err error) {
	if v.role != "judge" {
		return "", errors.New("only judge can traverse table")
	}

	// 读取table
	var table_dir string
	if v_task_id != "" {
		table_dir = filepath.Join(v.data_local_path, "verifier", task_id, v_task_id)
	} else {
		table_dir = filepath.Join(v.data_local_path, "executer", task_id)
	}

	// 查看是否符合表结构
	tables := new(DBVisitTask)
	err = utils.JsonToTable(table_dir, tables)
	if err != nil {
		return "", err
	}

	bytes, err := utils.TableToJson(tables)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (v *Data_api) GetTableCid(task_id, v_task_id string) (table string, err error) {
	if v.role != "judge" {
		return "", errors.New("only judge can get table cid")
	}

	// 读取table
	var table_dir string
	if v_task_id != "" {
		table_dir = filepath.Join(v.data_local_path, "verifier", task_id, v_task_id)
	} else {
		table_dir = filepath.Join(v.data_local_path, "executer", task_id)
	}

	cid, err := v.ipfs_api.GetFileCid(table_dir)
	if err != nil {
		return "", err
	}
	return cid, nil
}
