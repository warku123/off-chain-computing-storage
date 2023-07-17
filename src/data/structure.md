# Data部分存储文件夹结构
```sh
Data部分数据存储文件夹结构
data_local_path/data_ipns_name
├── db
├── excuter
│   ├── id1.json
│   ├── id2.json
│   └── id3.json
└── verfier
    └── id1
        ├── vid1.json
        ├── vid2.json
        └── vid3.json

executer里存储读表+写表
verifier里仅存储写表，用于验证者的写入

```