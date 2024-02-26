package main

import (
	"encoding/base64"
	"fmt"
	"github.com/libp2p/go-libp2p-core/crypto"
	"os"
)

func main() {
	// 生成一个新的密钥
	sk, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		panic(err)
	}

	// 编码密钥为base64
	skBytes, err := crypto.MarshalPrivateKey(sk)
	if err != nil {
		panic(err)
	}
	skB64 := base64.StdEncoding.EncodeToString(skBytes)

	// 输出swarm.key内容
	fmt.Printf("/key/swarm/psk/1.0.0/\n/base64/\n%s\n", skB64)

	// 将密钥写入文件
	file, err := os.Create("swarm.key")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString(fmt.Sprintf("/key/swarm/psk/1.0.0/\n/base64/\n%s\n", skB64))
}
