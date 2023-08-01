# 链下计算存储
## 运行
### 配置相关环境
运行以下命令以配置ipfs docker测试环境
```sh
docker pull warku123/ipfs-empty:v3 # 需要先连接dockerhub
docker run -d -p 5001:5001 warku123/ipfs-empty:v3
``` 
之后进入启动的docker环境
```sh
docker exec -it {container-hash} sh
```
在进入的docker环境中利用`ipfs key gen {name}`生成两个key,用于对应保存image及data

此处示例名字使用image和data

其中下方两个生成的IPNS哈希需要保存下来，其与密钥名字具有一一对应的关系，需要在API会话创建时与密钥名字一一对应输入
```sh
/ $ ipfs key gen image
k51qzi5uqu5dkfrh8bbp51n11hu78tzo6p8yds303qsk9vaalghy9kvyjh8cc0
/ $ ipfs key gen data
k51qzi5uqu5dk92lmmkmkahghprz5bejqsi32yr35qpp5sw81i2wbhe4822gi7
```
之后若需要进行功能测试，则断开公网连接
```sh
/ $ ipfs bootstrap rm --all
ipfs: Reading from /dev/stdin; send Ctrl-d to stop.
removed /dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN
removed /dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa
removed /dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb
removed /dnsaddr/bootstrap.libp2p.io/p2p/QmcZf59bWwK5XFi76CZX8cbJ4BhTzzA3gU1ZjYZcYW3dwt
removed /ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ
removed /ip4/104.131.131.82/udp/4001/quic/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ
```
最后`ctrl+D`退出docker环境，并运行以下命令重启docker，即可开始测试
```sh
docker restart {container-hash}
```
### 运行API
本地直接执行`go run ./main/main.go`

## API使用文档
https://documenter.getpostman.com/view/26790339/2s9XxvRtzj
