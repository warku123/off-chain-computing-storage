package image

import (
	"crypto/sha256"
	"fmt"

	"github.com/cbergoon/merkletree"
)

// 实现Interface的两个函数
func (v Image) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(v.Hash + v.Timestamp)); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func (v Image) Equals(other merkletree.Content) (bool, error) {
	return v.Hash == other.(Image).Hash && v.Timestamp == other.(Image).Timestamp, nil
}

func (v *Image_api) BuildTree() (err error) {
	var content []merkletree.Content
	if len((*v.image_table)[v.task_name]) == 0 {
		fmt.Println("No image in this task")
		return nil
	}

	for _, image := range (*v.image_table)[v.task_name] {
		content = append(content, image)
	}

	v.merkle_tree, err = merkletree.NewTree(content)
	if err != nil {
		return err
	}
	return nil
}

func (v *Image_api) GetRootHash() (hash string, err error) {
	hash = v.merkle_tree.Root.String()
	return hash, nil
}

func (v *Image_api) VerifyTree() (result bool, err error) {
	return v.merkle_tree.VerifyTree()
}

func (v *Image_api) VerifyContent(content Image) (result bool, err error) {
	return v.merkle_tree.VerifyContent(content)
}
