package image

import (
	"errors"
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
	if idx-offset > len((*v)[task_owner][task_name].Images) || idx-offset < 0 {
		return "", 0, errors.New("index out of range")
	}
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

func (v *ImageTable) GetOwnerImages(task_owner string) (result map[string]string, err error) {
	result = make(map[string]string)
	for task_name, tuple := range (*v)[task_owner] {
		if len(tuple.Images) == 0 {
			continue
		}
		result[task_name] = tuple.Images[0].Hash
	}
	return result, nil
}

func (v *ImageTable) GarbageCollection(task_owner string, task_name string, thershold uint64) error {
	// 从前往后找
	entry := (*v)[task_owner][task_name]
	i := 0
	for i = 0; i < len((*v)[task_owner][task_name].Images); i++ {
		if (*v)[task_owner][task_name].Images[i].Height <= thershold {
			entry.Offset += 1
		} else {
			break
		}
	}
	entry.Images = entry.Images[i:]
	(*v)[task_owner][task_name] = entry

	return nil
}
