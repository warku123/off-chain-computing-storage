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
		image.ImageWithIpnsName("k51qzi5uqu5did01y4bfh94mbd1olkqyyyj1hqhtrrqsxh97funiqyod9l2dx8"),
		image.ImageWithLocalPath("/Users/jojo/test/image/"),
	)
	if err != nil {
		t.Fatalf(err.Error())
	}
}
