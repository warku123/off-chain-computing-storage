package storage_off

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

// 写变量表，目前需要把所有value都放在内存中，并且只支持文本文件
type write_variable struct {
	Name   string   `json:"name"`
	Values []string `json:"values"` // 用于多次写的情况
}

// 读变量表，目前需要把所有value都放在内存中，并且只支持文本文件
type read_variable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// 单次计算任务的重复访问
type DBVisitTask struct {
	Read_table  []read_variable  `json:"read_table"`
	Write_table []write_variable `json:"write_table"`
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

// ipns发布镜像，返回ipns中的名字
func (v *ipfs_api) PublishDB(cid string) (string, error) {
	// 这边这个函数传的key的参数，不知道是key的name还是啥
	response, err := v.sh.PublishWithDetails(fmt.Sprintln("/ipfs/"+cid),
		v.data_key_name,
		24*time.Hour,
		24*time.Hour,
		true,
	)
	if err != nil {
		return "", err
	}

	return response.Name, nil
}

func (v *ipfs_api) InitDBVisit() (err error) {
	// 整个数据库持久化，速度有些慢，但是IPFS目前好像不支持数据库中部分文件持久化到本地
	err = v.FetchDB()
	if err != nil {
		return err
	}
	// 初始化DBVisitTask中的role
	// v.data_visit_task.role = v.role

	// 验证者，读变量表从IPFS转存到内存
	if v.role == "verifier" {
		read_table_path := fmt.Sprintf("%s/read_tables/%s_%d.json", v.data_local_path, v.snapshot_tag, v.image_idx)

		dataEncoded, err := ioutil.ReadFile(read_table_path)
		if err != nil {
			return err
		}
		json.Unmarshal(dataEncoded, &v.data_visit_task.Read_table)
	}
	return nil
}

func (v *ipfs_api) EndDBVisit() (err error) {
	// 首先Fetch整个DB
	err = v.FetchDB()
	if err != nil {
		return err
	}

	// 执行者时，存储读变量表
	if v.role == "executer" {
		read_table, err := json.Marshal(v.data_visit_task.Read_table)
		if err != nil {
			return err
		}

		read_table_path := fmt.Sprintf("%s/read_tables/%s_%d.json", v.data_local_path, v.snapshot_tag, v.image_idx)
		err = ioutil.WriteFile(read_table_path, read_table, 0666)
		if err != nil {
			return err
		}
	}

	// 存储写变量表
	read_table, err := json.Marshal(v.data_visit_task.Write_table)
	if err != nil {
		return err
	}

	write_table_path := fmt.Sprintf("%s/write_tables/%s_%d_%s.json", v.data_local_path, v.snapshot_tag, v.image_idx, v.role)
	err = ioutil.WriteFile(write_table_path, read_table, 0666)
	if err != nil {
		return err
	}

	cid, err := v.AddFolder(v.data_local_path)
	if err != nil {
		return err
	}

	_, err = v.PublishDB(cid)
	if err != nil {
		return err
	}

	return nil
}

// 输入变量的键以及对应版本，返回值
func (v *ipfs_api) ReadDB(key string, version int) (value string, err error) {
	// 变量在写表中
	for _, val := range v.data_visit_task.Write_table {
		if val.Name == key {
			return val.Values[len(val.Values)-1], nil
		}
	}

	// 变量在读表中
	for _, val := range v.data_visit_task.Read_table {
		if val.Name == key {
			return val.Value, nil
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

		v.data_visit_task.Read_table = append(v.data_visit_task.Read_table,
			read_variable{
				Name:  key,
				Value: string(data),
			})

		return string(data), nil
	}

	return "", errors.New(fmt.Sprintf("No file %s with version %s", key, strconv.Itoa(version)))
}

// 写入写表
func (v *ipfs_api) WriteDB(key string, value string) (err error) {
	// 变量已经存在
	var_exist := false
	for idx, val := range v.data_visit_task.Write_table {
		if val.Name == key {
			v.data_visit_task.Write_table[idx].Values = append(
				v.data_visit_task.Write_table[idx].Values, value)
			var_exist = true
			break
		}
	}
	// 变量不存在，新建变量
	if !var_exist {
		v.data_visit_task.Write_table = append(v.data_visit_task.Write_table, write_variable{
			Name:   key,
			Values: []string{value},
		})
	}
	return nil
}
