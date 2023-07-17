package data

import (
	"errors"
	"offstorage/json_op"
	"os"
	"path/filepath"
)

func (v *Data_api) DeleteTables() (err error) {
	if v.role != "judge" {
		return errors.New("only judge can delete tables")
	}

	exe_table_dir := filepath.Join(v.data_local_path, "executer", v.task_id)
	err = os.Remove(exe_table_dir)
	if err != nil {
		return err
	}

	ver_table_dir := filepath.Join(v.data_local_path, "verifier", v.task_id)
	ver_tables, err := filepath.Glob(ver_table_dir + "*")
	if err != nil {
		return err
	}
	for _, ver_table := range ver_tables {
		err = os.Remove(ver_table)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *Data_api) DataPersistance() (err error) {
	if v.role != "judge" {
		return errors.New("only judge can persist data")
	}

	for key, entry := range v.tables.Write_table {
		cid := entry[len(entry)-1]
		err = v.db.AddCid(key, cid)
		if err != nil {
			return err
		}
		err = v.db.ReduceWriteNum(key)
		if err != nil {
			return err
		}
	}

	// 上传table
	_, err = v.SyncDataToIPFS()
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
		new_data_list := make([]Data, 0)
		for i := 0; i < len(entry.Data_tuples); i++ {
			if entry.Data_tuples[i].Read_num > 0 {
				new_data_list = append(new_data_list, entry.Data_tuples[i])
			} else {
				entry.Gc_offset++
			}
		}
		entry.Data_tuples = new_data_list
		(*(v.db))[key] = entry
	}
	return nil
}

func (v *Data_api) InitDB() (err error) {
	if v.role != "judge" {
		return errors.New("only judge can init db")
	}

	// v.db = new(Data_table)
	err = json_op.GenEmptyTable(
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
