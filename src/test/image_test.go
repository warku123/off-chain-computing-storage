package test

import (
	"offstorage/image"
	"testing"
)

func TestImageInit(t *testing.T) {
	t.Log("Begin")
	_, err := image.NewImageShell(
		image.ImageWithHost("127.0.0.1"),
		image.ImageWithPort(5001),
		image.ImageWithTaskName("test"),
		image.ImageWithKeyName("test"),
		image.ImageWithIpnsName("k51qzi5uqu5dh3d6te4990h45jqr4mbh4e7om7ao88pvsfhi3qct15av2gwx5p"),
		image.ImageWithLocalPath("/Users/jojo/Documents/GitHub/off-chain-computing-storage/testfile"),
	)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestImageTableCreate(t *testing.T) {
	image_dir := "/Users/jojo/Documents/GitHub/off-chain-computing-storage/testfile/k51qzi5uqu5dh3d6te4990h45jqr4mbh4e7om7ao88pvsfhi3qct15av2gwx5p"
	image_api := make(image.ImageTable)
	// image_api.AddImageTuple("hash1", "123321", "test")
	// image_api.AddImageTuple("hash2", "123123", "test")
	// image_api.AddImageTuple("hash3", "123321123", "test2")
	err := image_api.SaveImageTable(image_dir)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestAddImage(t *testing.T) {
	image_api, err := image.NewImageShell(
		image.ImageWithHost("127.0.0.1"),
		image.ImageWithPort(5001),
		image.ImageWithOwnerName("test_owner"),
		image.ImageWithTaskName("test_task2"),
		image.ImageWithKeyName("test"),
		image.ImageWithIpnsName("k51qzi5uqu5dh3d6te4990h45jqr4mbh4e7om7ao88pvsfhi3qct15av2gwx5p"),
		image.ImageWithLocalPath("/Users/jojo/Documents/GitHub/off-chain-computing-storage/testfile"),
	)
	if err != nil {
		t.Fatalf(err.Error())
	}

	test_height := 10001
	cid, idx, err := image_api.AddImage(
		"/Users/jojo/Documents/GitHub/off-chain-computing-storage/testfile/test.txt",
		uint64(test_height),
	)
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Log("cid:" + cid)
	t.Logf("idx:%d", idx)

	height, err := image_api.GetImageByIdx(idx, "/Users/jojo/Documents/GitHub/off-chain-computing-storage/testfile/download.txt")
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Log(height, test_height)

	err = image_api.GetImageByCid(cid, "/Users/jojo/Documents/GitHub/off-chain-computing-storage/testfile/download2.txt")
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestTraverse(t *testing.T) {
	image_api, err := image.NewImageShell(
		image.ImageWithHost("127.0.0.1"),
		image.ImageWithPort(5001),
		image.ImageWithOwnerName("test_owner"),
		image.ImageWithTaskName("test_task"),
		image.ImageWithKeyName("test"),
		image.ImageWithIpnsName("k51qzi5uqu5dh3d6te4990h45jqr4mbh4e7om7ao88pvsfhi3qct15av2gwx5p"),
		image.ImageWithLocalPath("/Users/jojo/Documents/GitHub/off-chain-computing-storage/testfile"),
	)
	if err != nil {
		t.Fatalf(err.Error())
	}

	result, err := image_api.GetImageList()
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Log(result)
}
