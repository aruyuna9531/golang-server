# go

$ go run main.go

mytcp可以搭建一个tcp服务器，可以通过运行client.go测试tcp服务。

myhttp可以搭建一个http服务器，可以通过curl http://127.0.0.1:xxxx 测试http服务。

支持zookeeper。可以搭建一个小型微服务。

zookeeper需要安装[第三方包](https://github.com/samuel/go-zookeeper)的zookeeper库，并把里面的zk目录放到GOROOT里，才能编译通过。

支持连接MySQL。需要安装go的MySQL驱动：go get github.com/go-sql-driver/mysql

如果go get连不上，可以先执行export GOPROXY=https://goproxy.io