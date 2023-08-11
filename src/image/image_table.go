package image

import (
	"offstorage/utils"
)

type Image struct {
	Hash   string `json:"hash"`
	Height uint64 `json:"height"`
}

type ImageTable map[string][]Image

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

func (v *ImageTable) AddImageTuple(hash string, height uint64, task_name string) (idx int, err error) {
	image := Image{
		Hash:   hash,
		Height: height,
	}

	idx = len((*v)[task_name])
	(*v)[task_name] = append((*v)[task_name], image)

	return idx, nil
}

func (v *ImageTable) GetImageTuple(task_name string, idx int) (hash string, height uint64, err error) {
	image_tuple := (*v)[task_name][idx]

	return image_tuple.Hash, image_tuple.Height, nil
}

func (v *ImageTable) SaveImageTable(image_dir string) error {
	err := utils.SaveTable(image_dir, v)
	if err != nil {
		return err
	}
	return err
}
