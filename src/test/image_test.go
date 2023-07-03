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
		image.ImageWithChainName("test"),
		image.ImageWithKeyName("image"),
		image.ImageWithIpnsName("k51qzi5uqu5di645l8hd865kitoe5o29c2skixwpkgmemw24ffd924x54dan5a"),
		image.ImageWithLocalPath("/home/jzhang/ipfs_test/image"),
	)
	if err != nil {
		t.Fatalf(err.Error())
	}
}
