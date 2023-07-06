package json_op

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

func SaveJsonToFile(image_dir string, jsonBytes []byte) error {
	err := os.WriteFile(image_dir, jsonBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
