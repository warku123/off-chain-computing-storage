package data

import "errors"

type Data struct {
	value_hash string `json:"value_hash"`
	read_num   int    `json:"read_num"`
}

type Data_entry struct {
	data_tuples []Data `json:"data_tuples"`
	write_num   int    `json:"write_num"`
}

type Data_table map[string]Data_entry

func (v *Data_table) AddCid(name, hash string) {
	if entry, ok := (*v)[name]; ok {
		entry.data_tuples = append(entry.data_tuples, Data{
			value_hash: hash,
			read_num:   0,
		})
	} else {
		(*v)[name] = Data_entry{
			data_tuples: []Data{
				{
					value_hash: hash,
					read_num:   0,
				},
			},
			write_num: 0,
		}
	}
}

func (v *Data_table) GetCid(name string, version int) string {
	return (*v)[name].data_tuples[version].value_hash
}

func (v *Data_table) AddWriteNum(name string) {
	if entry, ok := (*v)[name]; ok {
		entry.write_num++
	} else {
		(*v)[name] = Data_entry{
			data_tuples: []Data{},
			write_num:   1,
		}
	}
}

func (v *Data_table) ReduceWriteNum(name string) error {
	if entry, ok := (*v)[name]; ok {
		entry.write_num--
	} else {
		return errors.New("ReduceWriteNum: no such data entry")
	}
	return nil
}

func (v *Data_table) GetDataVersionNum(name string) (int, error) {
	if _, ok := (*v)[name]; ok {
		return len((*v)[name].data_tuples), nil
	} else {
		return 0, errors.New("GetDataVersionNum: no such data entry")
	}
}

func (v *Data_table) AddReadNum(name string, version int) {
	(*v)[name].data_tuples[version].read_num++
}

func (v *Data_table) ReduceReadNum(name string, version int) {
	(*v)[name].data_tuples[version].read_num--
}
