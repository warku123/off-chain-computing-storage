package image

import (
	"offstorage/utils"
)

type Image struct {
	Hash   string `json:"hash"`
	Height uint64 `json:"height"`
}

type ImageTuple struct {
	Images []Image `json:"images"`
	Offset int     `json:"offset"`
}

type ImageTable map[string]map[string]ImageTuple // task_owner -> task_name -> image

func (v *ImageTable) GenEmptyImageTable(image_dir string) error {
	v = new(ImageTable)
	jsonBytes, err := utils.TableToJson(v)
	if err != nil {
		return err
	}

	err = utils.SaveJsonToFile(image_dir, jsonBytes)
	if err != nil {
		return err
	}
	return nil
}

func (v *ImageTable) AddImageTuple(hash string, height uint64, task_owner string, task_name string) (idx int, err error) {
	image := Image{
		Hash:   hash,
		Height: height,
	}

	table, exists := (*v)[task_owner]
	if !exists {
		table = make(map[string]ImageTuple)
		(*v)[task_owner] = table
	}

	tuple, exists := table[task_name]
	if !exists {
		tuple = ImageTuple{
			Images: make([]Image, 0),
			Offset: 0,
		}
	}
	idx = len(tuple.Images) + tuple.Offset

	tuple.Images = append(tuple.Images, image)
	// 复制回来
	table[task_name] = tuple
	(*v)[task_owner] = table

	return idx, nil
}

func (v *ImageTable) GetImageTuple(task_owner string, task_name string, idx int) (hash string, height uint64, err error) {
	offset := (*v)[task_owner][task_name].Offset
	image_tuple := (*v)[task_owner][task_name].Images[idx-offset]
	return image_tuple.Hash, image_tuple.Height, nil
}

func (v *ImageTable) SaveImageTable(image_dir string) error {
	err := utils.SaveTable(image_dir, v)
	if err != nil {
		return err
	}
	return err
}
