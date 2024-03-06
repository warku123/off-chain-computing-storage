package data

import "errors"

type Data struct {
	Value_hash string `json:"value_hash"`
	Read_num   int    `json:"read_num"`
}

type Data_entry struct {
	Data_tuples []Data `json:"data_tuples"`
	Write_num   int    `json:"write_num"`
	Gc_offset   int    `json:"gc_offset"` // 因为gc导致的index不对齐，需要一个偏移量调整
}

// 真正的要读取的index是version-offset
// 因此GetDataVersionNum时候是数组+偏移量

type Data_table map[string]Data_entry

func (v *Data_table) AddCid(name, hash string) (err error) {
	// 存在该key
	if entry, ok := (*v)[name]; ok {
		if entry.Write_num == 0 {
			// 因为修改写表时Writenum+1，所以不会出现0的情况
			return errors.New("add cid: write num is 0")
		} else if entry.Write_num == 1 &&
			len(entry.Data_tuples) != 0 &&
			entry.Data_tuples[len(entry.Data_tuples)-1].Read_num > 0 {
			// 只有一个写入任务，且最后一个数据没有读取任务，同时不是新建数据时
			// 直接修改最后一个版本数据
			entry.Data_tuples[len(entry.Data_tuples)-1].Value_hash = hash
		} else {
			// 其他情况需要新建一个版本数据
			entry.Data_tuples = append(entry.Data_tuples, Data{
				Value_hash: hash,
				Read_num:   0,
			})
		}
		(*v)[name] = entry
	} else {
		(*v)[name] = Data_entry{
			Data_tuples: []Data{
				{
					Value_hash: hash,
					Read_num:   0,
				},
			},
			Write_num: 0,
		}
	}
	return nil
}

func (v *Data_table) GetCid(name string, version int) (string, error) {
	offset := (*v)[name].Gc_offset
	if version-offset < 0 {
		return "", errors.New("GetCid: version-offset < 0")
	}
	return (*v)[name].Data_tuples[version-offset].Value_hash, nil
}

func (v *Data_table) AddWriteNum(name string) error {
	if entry, ok := (*v)[name]; ok {
		entry.Write_num++
		(*v)[name] = entry
	} else {
		(*v)[name] = Data_entry{
			Data_tuples: []Data{},
			Write_num:   1,
			Gc_offset:   0,
		}
	}
	return nil
}

func (v *Data_table) ReduceWriteNum(name string) error {
	if entry, ok := (*v)[name]; ok {
		if entry.Write_num <= 0 {
			return errors.New("ReduceWriteNum: write num <= 0")
		}
		entry.Write_num--
		(*v)[name] = entry
	} else {
		return errors.New("ReduceWriteNum: no such data entry")
	}
	return nil
}

// 某数据总共的版本数量
func (v *Data_table) GetDataVersionNum(name string) (int, error) {
	if _, ok := (*v)[name]; ok {
		return len((*v)[name].Data_tuples) + (*v)[name].Gc_offset, nil
	} else {
		return 0, errors.New("GetDataVersionNum: no such data entry")
	}
}

func (v *Data_table) AddReadNum(name string, version int) (err error) {
	offset := (*v)[name].Gc_offset
	(*v)[name].Data_tuples[version-offset].Read_num++
	return nil
}

func (v *Data_table) ReduceReadNum(name string, version int) (err error) {
	offset := (*v)[name].Gc_offset
	if (*v)[name].Data_tuples[version-offset].Read_num <= 0 {
		return errors.New("ReduceReadNum: read num <= 0")
	}
	(*v)[name].Data_tuples[version-offset].Read_num--
	return nil
}
