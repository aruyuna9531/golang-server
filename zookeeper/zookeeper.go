// 这里需要下载 https://github.com/samuel/go-zookeeper 的zookeeper库，并把里面的zk目录放到GOROOT里
// 目前先把demo搬过来

package zookeeper

import (
	"fmt"
	"time"
	"zk"
)

func LinkZookeeper() {
	fmt.Println("linking zookeeper... ")
	c, _, err := zk.Connect([]string{"127.0.0.1:2181"}, time.Second) // 这里要连集群时string数组里加
	if err != nil {
		fmt.Println("error, exiting(1): " + err.Error())
		panic(err)
	}
	for {
		children, stat, ch, err := c.ChildrenW("/")						// 要监听的节点路径，这个目录有节点变化时会通知到
		if err != nil {
			fmt.Println("error, exiting(2): " + err.Error())
			panic(err)
		}
		fmt.Printf("%+v %+v\n", children, stat)
		e := <-ch
		fmt.Printf("%+v\n", e)
	}
}
