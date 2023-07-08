package data

import "errors"

type Data struct {
	Value_hash string `json:"value_hash"`
	Read_num   int    `json:"read_num"`
}

type Data_entry struct {
	Data_tuples []Data `json:"data_tuples"`
	Write_num   int    `json:"write_num"`
}

type Data_table map[string]Data_entry

func (v *Data_table) AddCid(name, hash string) {
	if entry, ok := (*v)[name]; ok {
		entry.Data_tuples = append(entry.Data_tuples, Data{
			Value_hash: hash,
			Read_num:   0,
		})
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
}

func (v *Data_table) GetCid(name string, version int) string {
	return (*v)[name].Data_tuples[version].Value_hash
}

func (v *Data_table) AddWriteNum(name string) {
	if entry, ok := (*v)[name]; ok {
		entry.Write_num++
	} else {
		(*v)[name] = Data_entry{
			Data_tuples: []Data{},
			Write_num:   1,
		}
	}
}

func (v *Data_table) ReduceWriteNum(name string) error {
	if entry, ok := (*v)[name]; ok {
		entry.Write_num--
	} else {
		return errors.New("ReduceWriteNum: no such data entry")
	}
	return nil
}

func (v *Data_table) GetDataVersionNum(name string) (int, error) {
	if _, ok := (*v)[name]; ok {
		return len((*v)[name].Data_tuples), nil
	} else {
		return 0, errors.New("GetDataVersionNum: no such data entry")
	}
}

func (v *Data_table) AddReadNum(name string, version int) {
	(*v)[name].Data_tuples[version].Read_num++
}

func (v *Data_table) ReduceReadNum(name string, version int) {
	(*v)[name].Data_tuples[version].Read_num--
}
