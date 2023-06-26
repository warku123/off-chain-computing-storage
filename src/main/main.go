package main

import (
	"flag"
	"fmt"
	"offstorage/image"
	"os"
)

var host = flag.String("host", "127.0.0.1", "Host of IPFS")
var port = flag.Int("port", 5001, "Port of IPFS")
var role = flag.String("role", "verifier", "Role to visit IPFS")
var image_local_path = flag.String("ipath", "~/images", "Local path to save images")
var data_local_path = flag.String("dpath", "~/off_chain_data", "Local path to save data")
var image_key_name = flag.String("ikey", "test", "Key for image storage")
var data_key_name = flag.String("dkey", "test", "Key for data storage")
var tag = flag.String("tag", "default", "Chain name")

func main() {
	flag.Parse()
	ipfs, err := image.NewShell(
		image.ShellWithHost(*host),
		image.ShellWithPort(*port),
		image.ShellWithRole(*role),
		image.ShellWithImageKeyName(*image_key_name),
		image.ShellWithDataKeyName(*data_key_name),
		image.ShellWithTag(*tag),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	cid, err := ipfs.AddFile("./testfile/2023062409400900285981")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	fmt.Println(cid)
}
