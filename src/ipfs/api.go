package ipfs

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
)

// 添加字符串返回cid
func (v *Ipfs_api) AddString(value string) (cid string, err error) {
	cid, err = v.Sh.Add(strings.NewReader(value))
	fmt.Printf("added %s\n", cid)
	if err != nil {
		return "", err
	}
	return cid, nil
}

// 添加文件返回cid
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

// 添加文件夹返回cid
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

// 以cid寻址文件，下载到outdir
func (v *Ipfs_api) GetFile(cid string, outdir string) (err error) {
	err = v.Sh.Get(cid, outdir)
	return err
}

// ipns发布镜像，返回ipns中的名字，输入ipfs中cid和keyname
func (v *Ipfs_api) PublishFile(cid, keyname string) (ipnsname string, err error) {
	response, err := v.Sh.PublishWithDetails(
		filepath.Join("/ipfs/", cid),
		keyname,
		24*time.Hour,
		24*time.Hour,
		true,
	)
	if err != nil {
		return "", err
	}

	return response.Name, nil
}

func (v *Ipfs_api) GetFileCid(path string) (cid string, err error) {
	reader, err := os.Open(path)
	if err != nil {
		return "", err
	}

	cid, err = v.Sh.Add(reader, shell.OnlyHash(true))
	if err != nil {
		return "", err
	}
	return cid, nil
}
