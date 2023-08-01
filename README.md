# Off Chain Computing Storage
## 目标
利用IPFS实现一个链下计算存储
- 实现镜像和数据的分离存储
- 保证并行链下计算的数据正确访问
- 实现垃圾回收机制，保证数据库体积不会过大
## 项目结构
### deperecated_src
之前写的第一版代码，有一些功能和现在设计不符，因此弃用
### dockerfile
配置环境用的dockerfile，待更新
### src
链下计算代码部分
### testfile
用于测试的一些数据（已被gitignore屏蔽）