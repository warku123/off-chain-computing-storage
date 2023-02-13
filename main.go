package storage_off

import (
	"flag"
	"fmt"
	"os"
)

var host = flag.String("host", "127.0.0.1", "Host of IPFS")
var port = flag.Int("port", 5001, "Port of IPFS")
var role = flag.String("role", "verifier", "Role to visit IPFS")
var image_local_path = flag.String("ipath", "/images", "Local path to save images")
var data_local_path = flag.String("dpath", "/off_chain_data", "Local path to save data")
var image_key_name = flag.String("ikey", "", "Key for image storage")
var data_key_name = flag.String("dkey", "", "Key for data storage")
var tag = flag.String("tag", "", "Chain name")

func main() {
	flag.Parse()
	ipfs, err := NewShell(
		ShellWithHost(*host),
		ShellWithPort(*port),
		ShellWithRole(*role),
		ShellWithImageKeyName(*image_key_name),
		ShellWithDataKeyName(*data_key_name),
		ShellWithTag(*tag),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	err = Interactive_shell(ipfs)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
