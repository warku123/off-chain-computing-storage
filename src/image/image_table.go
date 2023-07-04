package image

import (
	"encoding/json"
	"os"
)

type Image struct {
	Hash      string `json:"hash"`
	Timestamp string `json:"timestamp"`
}

type ImageTable map[string][]Image

func (v *ImageTable) JsonToImageTable(image_dir string) error {
	if v == nil {
		v = new(ImageTable)
	}

	content, err := os.ReadFile(image_dir)
	if err != nil {
		return err
	}
	// fmt.Println("content" + string(content))
	err = json.Unmarshal(content, v)
	if err != nil {
		return err
	}
	return nil
}

func (v *ImageTable) ImageTableToJson() ([]byte, error) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return jsonBytes, nil
}

func (v *ImageTable) SaveJsonToFile(image_dir string, jsonBytes []byte) error {
	err := os.WriteFile(image_dir, jsonBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (v *ImageTable) GenEmptyImageTable(image_dir string) error {
	v = new(ImageTable)
	jsonBytes, err := v.ImageTableToJson()
	if err != nil {
		return err
	}

	err = v.SaveJsonToFile(image_dir, jsonBytes)
	if err != nil {
		return err
	}
	return nil
}

func (v *ImageTable) AddImageTuple(hash, timestamp, task_name string) (idx int, err error) {
	image := Image{
		Hash:      hash,
		Timestamp: timestamp,
	}

	idx = len((*v)[task_name])
	(*v)[task_name] = append((*v)[task_name], image)

	return idx, nil
}

func (v *ImageTable) GetImageTuple(task_name string, idx int) (hash, timestamp string, err error) {
	image_tuple := (*v)[task_name][idx]

	return image_tuple.Hash, image_tuple.Timestamp, nil
}

func (v *ImageTable) SaveImageTable(image_dir string) error {
	jsonBytes, err := v.ImageTableToJson()
	if err != nil {
		return err
	}
	err = v.SaveJsonToFile(image_dir, jsonBytes)
	if err != nil {
		return err
	}
	return err
}
