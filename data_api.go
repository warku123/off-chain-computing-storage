package storage_off

import "fmt"

// 将持久化数据库下载到本地
func (v *ipfs_api) FetchDB() error {
	ipns_path := fmt.Sprintf("/ipns/%s", v.data_ipns_id)
	err := v.sh.Get(ipns_path, v.data_local_path)
	if err != nil {
		return err
	}
	return nil
}
