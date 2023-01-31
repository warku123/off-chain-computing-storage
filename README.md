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