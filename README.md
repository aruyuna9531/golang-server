# SimpleServer

golang 基本服务器代码

如果运行失败，可以输入
```shell
go env -w GOPROXY=https://goproxy.cn
```

---------------

目前的更新进度：

能正常发收包（没做超大包校验和分包粘包，超大包搞崩服务器这种情况要避免一下 TODO）

可使用protobuf在C-S之间进行通信