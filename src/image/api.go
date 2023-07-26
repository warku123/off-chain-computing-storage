package image

import (
	"bytes"
	"path/filepath"
)

func (v *Image_api) PublishImageTable() (err error) {
	image_dir := filepath.Join(v.image_local_path, v.image_ipns_name)
	image_table_cid, err := v.ipfs_api.AddFile(image_dir)
	if err != nil {
		return err
	}

	_, err = v.ipfs_api.PublishFile(image_table_cid, v.image_key_name)
	if err != nil {
		return err
	}
	return nil
}

func (v *Image_api) AddImage(image_path, timestamp string) (cid string, idx int, err error) {
	cid, err = v.ipfs_api.AddFile(image_path)
	if err != nil {
		return "", -1, err
	}

	idx, err = v.image_table.AddImageTuple(cid, timestamp, v.task_name)
	if err != nil {
		return "", -1, err
	}

	image_table_path := filepath.Join(v.image_local_path, v.image_ipns_name)
	err = v.image_table.SaveImageTable(image_table_path)
	if err != nil {
		return "", -1, err
	}

	err = v.PublishImageTable()
	if err != nil {
		return "", -1, err
	}

	// build merkle tree
	err = v.BuildTree()
	if err != nil {
		return "", -1, err
	}

	return cid, idx, nil
}

func (v *Image_api) GetImageByIdx(idx int, outdir string) (timestamp string, err error) {
	cid, timestamp, err := v.image_table.GetImageTuple(v.task_name, idx)
	if err != nil {
		return "", err
	}

	err = v.ipfs_api.GetFile(cid, outdir)
	if err != nil {
		return "", err
	}
	return timestamp, nil
}

func (v *Image_api) GetImageByCid(cid, outdir string) (err error) {
	err = v.ipfs_api.GetFile(cid, outdir)
	if err != nil {
		return err
	}
	return nil
}

func (v *Image_api) CatImageByCid(cid string) (image string, err error) {
	content, err := v.ipfs_api.Sh.Cat(cid)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(content)
	finalStr := buf.String()

	return finalStr, nil
}

func (v *Image_api) CatImageByIdx(idx int) (image string, timestamp string, err error) {
	cid, timestamp, err := v.image_table.GetImageTuple(v.task_name, idx)
	if err != nil {
		return "", "", err
	}

	content, err := v.CatImageByCid(cid)
	if err != nil {
		return "", "", err
	}
	return content, timestamp, nil
}
