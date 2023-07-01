package ipfs

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

// 三个add函数，添加字符串、文件或文件夹，返回cid
func (v *Ipfs_api) AddString(value string) (cid string, err error) {
	cid, err = v.Sh.Add(strings.NewReader(value))
	fmt.Printf("added %s\n", cid)
	if err != nil {
		return "", err
	}
	return cid, nil
}

func (v *Ipfs_api) AddFile(path string) (cid string, err error) {
	reader, err := os.Open(path)
	if err != nil {
		return "", err
	}

	cid, err = v.Sh.Add(reader)
	fmt.Printf("added %s\n", cid)
	if err != nil {
		return "", err
	}
	return cid, nil
}

func (v *Ipfs_api) AddFolder(path string) (cid string, err error) {
	cid, err = v.Sh.AddDir(path)
	fmt.Printf("added %s\n", cid)
	if err != nil {
		return "", err
	}
	return cid, nil
}

// 以cid寻址文件，返回文件内容string
func (v *Ipfs_api) ReadFile(cid string) (string, error) {
	content, err := v.Sh.Cat(cid)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(content)
	finalStr := buf.String()
	return finalStr, nil
}

func (v *Ipfs_api) GetFile(cid string, outdir string) (err error) {
	err = v.Sh.Get(cid, outdir)
	return err
}
