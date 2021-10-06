# c++

目前服务器支持监听一个TCP和一个HTTP

测试TCP服务器用的客户端程序可以去golang目录go run client.go

测试HTTP服务器，直接找个地方curl即可

server进程启动后，log目录下a.log是服务器运行日志

需要安装的第三方库见references目录

先把全部第三方库make install，再执行本目录的make

make完在bin目录下执行server文件即可
