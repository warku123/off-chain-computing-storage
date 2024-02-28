#!/bin/sh

# 导入IPNS键
ipfs key import image /image.key
ipfs key import data /data.key

# 启动IPFS守护进程
exec ipfs daemon --migrate=true