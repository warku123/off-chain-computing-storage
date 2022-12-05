package storage_off

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
)

func (v *ipfs_api) AddString(value string) (string, error) {
	cid, err := v.sh.Add(strings.NewReader(value))
	fmt.Printf("added %s\n", cid)
	if err != nil {
		return "", err
	}
	return cid, nil
}

func (v *ipfs_api) AddFile(path string) (string, error) {
	reader, err := os.Open(path)
	if err != nil {
		return "", err
	}

	cid, err := v.sh.Add(reader)
	fmt.Printf("added %s\n", cid)
	if err != nil {
		return "", err
	}
	return cid, nil
}

func (v *ipfs_api) AddFolder(path string) (string, error) {
	cid, err := v.sh.AddDir(path)
	fmt.Printf("added %s\n", cid)
	if err != nil {
		return "", err
	}
	return cid, nil
}

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

func (v *ipfs_api) PublishFile(cid string) (string, error) {
	key_name := v.snapshot_tag
	response, err := v.sh.PublishWithDetails(fmt.Sprintln("/ipfs/"+cid),
		key_name,
		24*time.Hour,
		24*time.Hour,
		true,
	)
	if err != nil {
		return "", err
	}

	return response.Name, nil
}

func (v *ipfs_api) NewImage(image string) (image_cid string, idx int64, err error) {
	image_cid, err = v.AddString(image)
	if err != nil {
		return "", -1, err
	}

	ipns_path := fmt.Sprintf("/ipns/%s", v.ipns_id)
	table_dest_dir := fmt.Sprintf("%s/snapshot.txt", v.snapshot_tag)
	err = v.sh.Get(ipns_path, table_dest_dir)
	if err != nil {
		return "", -1, err
	}

	file, err := os.Open(table_dest_dir)
	if err != nil {
		return "", -1, err
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

	idx = int64(count) + 1
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

	v.ipns_id, err = v.PublishFile(file_cid)
	if err != nil {
		return "", -1, err
	}

	return image_cid, idx, nil
}

func (v *ipfs_api) SearchImageByIdx(idx int64) (image string, err error) {
	ipns_path := fmt.Sprintf("/ipns/%s", v.ipns_id)
	table_dest_dir := fmt.Sprintf("%s/snapshot.txt", v.snapshot_tag)
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
