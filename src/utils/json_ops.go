package utils

import (
	"encoding/json"
	"os"
)

func JsonToTable(image_dir string, v any) error {
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

func TableToJson(v any) ([]byte, error) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return jsonBytes, nil
}

func SaveJsonToFile(dir string, jsonBytes []byte) error {
	// 给文件加个锁
	file, err := LockFile(dir)
	if err != nil {
		return err
	}
	defer UnlockFile(file)
	err = os.WriteFile(dir, jsonBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GenEmptyTable(dir string) error {
	err := os.WriteFile(dir, []byte("{}"), 0644)
	if err != nil {
		return err
	}
	return nil
}

func SaveTable(dir string, v any) error {
	jsonBytes, err := TableToJson(v)
	if err != nil {
		return err
	}
	err = SaveJsonToFile(dir, jsonBytes)
	if err != nil {
		return err
	}
	return nil
}
