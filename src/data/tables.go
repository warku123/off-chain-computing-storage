package data

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
