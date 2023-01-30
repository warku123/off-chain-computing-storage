package storage_off

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
)

// 写变量表，目前需要把所有value都放在内存中，并且只支持文本文件
type write_variable struct {
	name   string
	values []string // 用于多次写的情况
}

// 读变量表，目前需要把所有value都放在内存中，并且只支持文本文件
type read_variable struct {
	name  string
	value string
}

// 单次计算任务的重复访问
type DBVisitTask struct {
	// role        string //来自计算者or验证者的任务，好像可以直接用ipfs_api中的role代替
	read_table  []read_variable
	write_table []write_variable
}

// 将持久化数据库下载到本地
func (v *ipfs_api) FetchDB() error {
	ipns_path := fmt.Sprintf("/ipns/%s", v.data_ipns_id)
	err := v.sh.Get(ipns_path, v.data_local_path)
	if err != nil {
		return err
	}
	return nil
}

// 将单个文件下载到本地
func (v *ipfs_api) FetchFile(key string, version int) error {
	ipns_path := fmt.Sprintf("/ipns/%s/%s_%s.txt", v.data_ipns_id, key, strconv.Itoa(version))
	local_path := fmt.Sprintf("%s/%s_%s.txt", v.data_local_path, key, strconv.Itoa(version))
	err := v.sh.Get(ipns_path, local_path)
	if err != nil {
		return err
	}
	return nil
}

func (v *ipfs_api) InitDBVisit() (err error) {
	// 貌似不用整个持久化，速度太慢
	// err = v.FetchDB()
	// if err != nil {
	// 	return err
	// }
	// 初始化DBVisitTask中的role
	// v.data_visit_task.role = v.role
	return nil
}

// 输入变量的键以及对应版本，返回值
func (v *ipfs_api) ReadDB(key string, version int) (value string, err error) {

	// 变量在写表中
	for _, val := range v.data_visit_task.write_table {
		if val.name == key {
			return val.values[len(val.values)-1], nil
		}
	}

	// 变量在读表中
	for _, val := range v.data_visit_task.read_table {
		if val.name == key {
			return val.value, nil
		}
	}

	// 变量只在db中，仅限执行者
	if v.role == "executer" {
		err = v.FetchFile(key, version)
		if err != nil {
			return "", nil
		}
		local_path := fmt.Sprintf("%s/%s_%s.txt", v.data_local_path, key, strconv.Itoa(version))
		data, err := ioutil.ReadFile(local_path)
		if err != nil {
			fmt.Println(err)
			return "", err
		}

		v.data_visit_task.read_table = append(v.data_visit_task.read_table,
			read_variable{
				name:  key,
				value: string(data),
			})

		return string(data), nil
	}

	return "", errors.New(fmt.Sprintf("No file %s with version %s", key, strconv.Itoa(version)))
}
