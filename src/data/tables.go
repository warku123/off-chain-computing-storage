package data

// 写变量表，目前需要把所有value都放在内存中
type write_variable []string

// 读变量表，目前需要把所有value都放在内存中
type read_variable struct {
	Value        string `json:"value"`
	Read_version int    `json:"read_version"`
}

// 单次计算任务的重复访问
type DBVisitTask struct {
	Read_table  map[string]read_variable  `json:"read_table"`
	Write_table map[string]write_variable `json:"write_table"`
}

func (v *DBVisitTask) AddReadTuple(name, value string, version int) {
	v.Read_table[name] = read_variable{
		Value:        value,
		Read_version: version,
	}
}

func (v *DBVisitTask) AddWriteTuple(name, value string) {
	v.Write_table[name] = append(v.Write_table[name], value)
}
