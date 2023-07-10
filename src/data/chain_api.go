package data

import (
	"os"
	"path/filepath"
)

func (v *Data_api) DeleteTables() (err error) {
	exe_table_dir := filepath.Join(v.data_local_path, v.data_ipns_name, "executer", v.task_id)
	err = os.Remove(exe_table_dir)
	if err != nil {
		return err
	}

	ver_table_dir := filepath.Join(v.data_local_path, v.data_ipns_name, "verifier", v.task_id)
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
	return nil
}

func (v *Data_api) ReduceReadNum() (err error) {
	for key, entry := range v.tables.Read_table {
		err = v.db.ReduceReadNum(key, entry.Read_version)
		if err != nil {
			return err
		}
	}
	return nil
}
