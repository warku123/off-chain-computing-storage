# off-chain-computing-storage
off chain computing storage build by IPFS

## 构成
- image 镜像部分存储
- data 数据部分存储

## 初始化
初始化用`NewShell`函数，并且必须带有一个以斜杠为开始的路径，用于存储索引镜像表

## 可能有的bug
- [ ] 当前访问的image索引由索引-cid表的条目数量提供，若两个任务并发访问，可能产生一样的image_id
- [ ] 每次对数据库进行修改的时候没有加锁，可能导致高并发时候出现问题

## 基本使用Guide
### `ipfs_init.go`
该部分函数用于初始化ipfs_api结构体，以构建对IPFS存储的访问
#### `NewShell`
```
初始化一个ipfs_api结构体
Input: 不定个数的初始化函数modIPFSAPI
Return: ipfs_api结构体
```
#### `ShellWithxxx`
```
给结构体安全传参
Input: 参数
Return: 初始化函数modIPFSAPI
```
### `data_api.go`
访问数据存储相关函数
#### `ReadDB`
```
读数据
Input: 要访问的键名key以及版本version
Return: 访问的值value
```
#### `WriteDB`
```
写数据
Input: 要写入的键名key以及值value
```
#### `DataPersistence`
```
持久化写表中内容
Input: 数据版本号列表version_list，需要与read_list中的变量一一对应
```
### `image_api.go`
访问镜像存储相关函数
#### `NewImage`
```
新建镜像
Input: 镜像内容image
Return: 存储的哈希cid，索引idx
```
#### `SearchImageByIdx`
```
根据索引寻找镜像
Input: 镜像索引idx
Return: 镜像内容
```
#### `SearchImageByIdx`
```
根据哈希寻找镜像
Input: 镜像索引cid
Return: 镜像内容
```