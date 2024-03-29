package image

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
)

// 三个add函数，添加字符串、文件或文件夹，返回cid
func (v *ipfs_api) AddString(value string) (cid string, err error) {
	cid, err = v.sh.Add(strings.NewReader(value))
	fmt.Printf("added %s\n", cid)
	if err != nil {
		return "", err
	}
	return cid, nil
}

func (v *ipfs_api) AddFile(path string) (cid string, err error) {
	reader, err := os.Open(path)
	if err != nil {
		return "", err
	}

	cid, err = v.sh.Add(reader)
	fmt.Printf("added %s\n", cid)
	if err != nil {
		return "", err
	}
	return cid, nil
}

func (v *ipfs_api) AddFolder(path string) (cid string, err error) {
	cid, err = v.sh.AddDir(path)
	fmt.Printf("added %s\n", cid)
	if err != nil {
		return "", err
	}
	return cid, nil
}

// 在IPFS上构建一个Image存储系统，有bug待改，未使用
func (v ipfs_api) BuildImage() (string, string, error) {
	dest_dir := v.image_local_path
	dir_stat, err := os.Stat(dest_dir)
	if err != nil {
		return "", "", err
	}
	dir_exist := dir_stat.IsDir()
	if !dir_exist {
		err := os.MkdirAll(dest_dir, os.ModePerm)
		if err != nil {
			return "", "", err
		}
	}

	table_dest_dir := fmt.Sprintf("%s/snapshot.txt", v.snapshot_tag)
	_, err = os.Stat(table_dest_dir)
	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(table_dest_dir)
			if err != nil {
				return "", "", err
			}
		}
	}

	image_key, err := v.genIPNSkey()
	if err != nil {
		return "", "", err
	}

	// image_ipns_cid 对应的是快照索引表
	cid, err := v.AddFile(table_dest_dir)
	if err != nil {
		return "", "", err
	}

	v.image_ipns_id, err = v.PublishImage(cid)
	if err != nil {
		return "", "", err
	}

	// 这边key.Id想返回key的内容，但是不知道key里面是不是，待测试
	return image_key.Id, v.image_ipns_id, err
}

// 获得当前的image index
func (v *ipfs_api) GetImageIdx(table_dest_dir string) (index int64, err error) {
	file, err := os.Open(table_dest_dir)
	if err != nil {
		return -1, err
	}

	fd := bufio.NewReader(file)
	count := 0
	for {
		_, err := fd.ReadString('\n')
		if err != nil {
			break
		}
		count++
	}
	index = int64(count) + 1

	return index, nil
}

// 初始化image存储，也就是下载一个索引_cid映射表
func (v *ipfs_api) InitImage() error {
	dest_dir := v.image_local_path
	table_ipns_path := fmt.Sprintf("/ipns/%s", v.image_ipns_id)

	// 下载远端的所有image索引-cid表到本地
	err := v.sh.Get(table_ipns_path, dest_dir)
	if err != nil {
		return err
	}

	// 查看远端有没有该snapshot_tag的文件夹
	folder_dest_dir := dest_dir + "/" + v.snapshot_tag
	dir_stat, err := os.Stat(folder_dest_dir)
	if err != nil {
		return err
	}
	// 若没有，创建文件夹&创建表
	dir_exist := dir_stat.IsDir()
	if !dir_exist {
		err := os.MkdirAll(dest_dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	table_dest_dir := fmt.Sprintf("%s/snapshot.txt", folder_dest_dir)
	_, err = os.Stat(table_dest_dir)
	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(table_dest_dir)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// 以cid寻址文件，返回文件内容string
func (v *ipfs_api) ReadFile(cid string) (string, error) {
	content, err := v.sh.Cat(cid)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(content)
	finalStr := buf.String()
	return finalStr, nil
}

// ipns发布镜像，返回ipns中的名字
func (v *ipfs_api) PublishImage(cid string) (string, error) {
	// 这边这个函数传的key的参数，不知道是key的name还是啥
	response, err := v.sh.PublishWithDetails(fmt.Sprintln("/ipfs/"+cid),
		v.image_key_name,
		24*time.Hour,
		24*time.Hour,
		true,
	)
	if err != nil {
		return "", err
	}

	return response.Name, nil
}

// 存储Image的过程
func (v *ipfs_api) NewImage(image string) (image_cid string, idx int64, err error) {
	image_cid, err = v.AddString(image)
	if err != nil {
		return "", -1, err
	}

	ipns_path := fmt.Sprintf("/ipns/%s", v.image_ipns_id)
	table_dest_dir := fmt.Sprintf("%s/%s/snapshot.txt", v.image_local_path, v.snapshot_tag)

	// 先将本地的快照索引表更新
	err = v.sh.Get(ipns_path, table_dest_dir)
	if err != nil {
		return "", -1, err
	}
	// 修改快照索引表
	idx, err = v.GetImageIdx(table_dest_dir)
	if err != nil {
		return "", -1, err
	}

	file, err := os.Open(table_dest_dir)
	if err != nil {
		return "", -1, err
	}
	content := fmt.Sprintf("%s\r\n", image_cid)
	write := bufio.NewWriter(file)
	write.WriteString(content)
	//Flush将缓存的文件真正写入到文件中
	write.Flush()

	file.Close()

	file_cid, err := v.AddFile(table_dest_dir)
	if err != nil {
		return "", -1, err
	}

	v.image_ipns_id, err = v.PublishImage(file_cid)
	if err != nil {
		return "", -1, err
	}

	return image_cid, idx, nil
}

func (v *ipfs_api) SearchImageByIdx(idx int64) (image string, err error) {
	ipns_path := fmt.Sprintf("/ipns/%s", v.image_ipns_id)
	table_dest_dir := fmt.Sprintf("%s/%s/snapshot.txt", v.image_local_path, v.snapshot_tag)

	err = v.sh.Get(ipns_path, table_dest_dir)
	if err != nil {
		return "", err
	}

	file, err := os.Open(table_dest_dir)
	if err != nil {
		return "", err
	}

	var image_cid string
	fileScanner := bufio.NewScanner(file)
	lineCount := int64(1)
	for fileScanner.Scan() {
		if lineCount == idx {
			image_cid = fileScanner.Text()
		}
		lineCount++
	}

	image, err = v.ReadFile(image_cid)
	if err != nil {
		return "", err
	}
	return image, nil
}

func (v *ipfs_api) SearchImageByCid(cid string) (string, error) {
	content, err := v.sh.Cat(cid)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(content)
	image := buf.String()
	return image, nil
}
