# c++

来自chenshuo大佬的demo，加了一点点自己的改装

从效果来看，只比大佬demo多了一个http服务，其他没有。

目前服务器支持监听一个TCP和一个HTTP

需要安装[muduo](https://github.com/chenshuo/muduo)

注意muduo源码要放在/usr/local/include、muduo动态lib要放在/usr/local/lib

然后make就行了，makefile我已经写好了

make完在bin目录下执行server文件即可
