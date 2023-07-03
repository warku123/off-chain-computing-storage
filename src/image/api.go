package image

import "fmt"

func (v *Image_api) PublishImageTable() (err error) {
	image_dir := fmt.Sprintf("%s/%s", v.image_local_path, v.image_ipns_name)
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

func (v *Image_api) AddImage(image_path, timestamp string) (idx int, err error) {
	cid, err := v.ipfs_api.AddFile(image_path)
	if err != nil {
		return -1, err
	}

	idx, err = v.image_table.AddImageTuple(cid, timestamp, v.chain_name)
	if err != nil {
		return -1, err
	}

	err = v.image_table.SaveImageTable(image_path)
	if err != nil {
		return -1, err
	}

	err = v.PublishImageTable()
	if err != nil {
		return -1, err
	}

	return idx, nil
}

func (v *Image_api) GetImageByIdx(idx int, outdir string) (timestamp string, err error) {
	cid, timestamp, err := v.image_table.GetImageTuple(v.chain_name, idx)
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
