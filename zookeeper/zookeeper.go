package zookeeper

import (
	"fmt"
	"time"
	"github.com/go-zookeeper/zk"
	"strconv"
	"encoding/json"
)

type ZkConf struct {
	Addr		[]string	`json:"addr"`
	ServerId	int			`json:"svr_id"`
	HttpListen	int			`json:"http"`
	TcpListen	int			`json:"tcp"`
}

// e的数据类型：zk.Event，结构如下
// type Event struct {
//	Type   EventType
// (ExistsW)：观察节点本身，EventNodeDataChanged/EventNodeDeleted/EventNodeCreated
// (ChildrenW)：观察节点的直接子节点，EventNodeChildrenChanged/EventNodeDeleted（如果自身被删除）
// (GetW)：观察数据，EventNodeDataChanged/EventNodeDeleted（如果自身被删除）
// 注意，如果是监听ChildrenW和GetW的时候znode不存在或者被删除，那么监听会失败，并且不会重试，即使后面create了也不会有反应（只有ExistsW会响应create）
//	State  State
//	Path   string   //响应事件的节点路径
//	Err    error
//	Server string   // For connection events
//}

var zkInit bool = false
var c *zk.Conn

func createZkPath(server_id int) string {
	return "/gosvr/" + strconv.Itoa(server_id)
}

// 监听一个znode的子节点变化
func listenChildNode(path string) {
	if !zkInit {
		fmt.Println("zookeeper does not init, return")
		return
	}
	for {
		children, _, ch, err := c.ChildrenW(path)
		// 要监听的节点路径，这个目录有节点变化时会通知到。
		// 返回参数1 children:这个路径当前的所有子节点列表
		// 返回参数2 stat：一些变化数据，见代码末尾①。stat在ExistW获取即可，同一个事件监听器返回的stat都是一样的，避免过多打印
		// 返回参数3 ch：一个channel，监听节点存在事件变化时会传入数据
		// 返回参数4 err：抛出错误
		// 这行不能放到循环外，否则会出现死循环（StateDisconnected）
	
		if err != nil {
			fmt.Println("error, " + err.Error())
			return
		}
		
		fmt.Printf("children list：%+v\n", children)
		e := <-ch

		switch e.Type {
			// ChildrenW不监听Created
		case zk.EventNodeDeleted:
			// 节点自己被删除了
			fmt.Printf("znode %s self deleted, by ChildrenW(exit peacefully)\n", path)
			return
		case zk.EventNodeChildrenChanged:
			// 子节点发生了变化
			fmt.Printf("znode %s children change detected, by ChildrenW\n", path)
		case zk.EventNodeDataChanged:
			// 节点值发生了变化（无法获得具体值）
			fmt.Printf("znode %s data change detected, by ChildrenW\n", path)
		case zk.EventNotWatching:
			// 节点失去观察（interrupt时会发生）
			fmt.Printf("znode %s data not watching, by ExistW(exit peacefully)\n", path)
			return
		default:
			// 未知事件
			fmt.Printf("znode %s unknown event detected, by ChildrenW, event id = %v\n", path, e.Type)
		}
	}
}

// 监听一个znode的值变化
func listenSelfValue(path string) {
	if !zkInit {
		fmt.Println("zookeeper does not init, return")
		return
	}
	for {
		value, _, ch, err := c.GetW(path)
		// 要监听的节点路径，这个目录有节点变化时会通知到。
		// 返回参数1 children:这个路径当前的所有子节点列表
		// 返回参数2 stat：一些变化数据，见代码末尾①
		// 返回参数3 ch：一个channel，监听节点存在事件变化时会传入数据
		// 返回参数4 err：抛出错误
		// 这行不能放到循环外，否则会出现死循环（StateDisconnected）
	
		if err != nil {
			fmt.Println("error, " + err.Error())
			return
		}
		
		fmt.Printf("new value is：%+v\n", string(value))
		e := <-ch

		switch e.Type {
			// GetW不监听Created和ChildrenChange
		case zk.EventNodeDeleted:
			// 节点自己被删除了
			fmt.Printf("znode %s self deleted, by GetW(exit peacefully)\n", path)
			return
		case zk.EventNodeDataChanged:
			// 节点值发生了变化
			fmt.Printf("znode %s children change detected, by GetW, new value is：%s\n", path, string(value))
		case zk.EventNotWatching:
			// 节点失去观察（interrupt时会发生）
			fmt.Printf("znode %s data not watching, by ExistW(exit peacefully)\n", path)
			return
		default:
			// 未知事件
			fmt.Printf("znode %s unknown event detected, by GetW, event id = %v\n", path, e.Type)
		}
	}
}

// 监听一个znode
func listenNode(path string) {
	if !zkInit {
		fmt.Println("zookeeper does not init, return")
		return
	}

	// 启动即存在时初始化其他watcher
	exist, _, _, _ := c.ExistsW(path)
	if exist {
		fmt.Printf("znode %s existing, by ExistW, add other watcher\n", path)
		go listenChildNode(path)
		go listenSelfValue(path)
	}

	for {
		exist, stat, ch, err := c.ExistsW(path)
		// 要监听的节点路径，这个目录有节点变化时会通知到。
		// 返回参数1 exist:节点是否存在
		// 返回参数2 stat：一些变化数据，见代码末尾①
		// 返回参数3 ch：一个channel，监听节点存在事件变化时会传入数据
		// 返回参数4 err：抛出错误
		// 这行不能放到循环外，否则会出现死循环（StateDisconnected）
	
		if err != nil {
			fmt.Println("error, " + err.Error())
			return
		}
		
		fmt.Printf("exist：%+v，stat：%+v\n", exist, stat)
		e := <-ch

		switch e.Type {
		case zk.EventNodeCreated:
			// 节点被创建
			fmt.Printf("znode %s created, by ExistW\n", path)
			go listenChildNode(path)
			go listenSelfValue(path)
		case zk.EventNodeDeleted:
			// 节点被删除
			fmt.Printf("znode %s deleted, by ExistW\n", path)
		case zk.EventNodeDataChanged:
			// 节点值修改（但是ExistW的watcher不会返回修改后的值，要通过GetW获取）
			fmt.Printf("znode %s data change detected, by ExistW\n", path)
		case zk.EventNotWatching:
			// 节点失去观察（interrupt时会发生）
			fmt.Printf("znode %s data not watching, by ExistW(exit peacefully)\n", path)
			return
		default:
			// 未知事件
			fmt.Printf("zonde %s receive unknown event, by ExistW, event id = %v\n", path, e.Type)
		}
	}
}

func DeleteSelfZnode() {
	c.Close()	//断开连接时所有创建的临时节点会被立即删除
}

func LinkZookeeper(conf *ZkConf) {
	fmt.Println("linking zookeeper... ")
	var err error
	c, _, err = zk.Connect(conf.Addr, time.Second) // 这里要连集群时string数组里加
	if err != nil {
		fmt.Println("error, exiting(1): " + err.Error())
		return
	}

	zkInit = true

	zkValueStr, err := json.Marshal(conf)
	if err != nil {
		fmt.Println("error, exiting(2): " + err.Error())
		return
	}

	// 把自己注册到zookeeper。zk.FlagEphemeral 创建临时节点，服务器进程结束后会删掉节点。但是有一定延时，如果进程关掉立马重启是会出现“节点已存在”的提示的
	if _, err = c.Create(createZkPath(conf.ServerId), zkValueStr, zk.FlagEphemeral, zk.WorldACL(zk.PermAll)); err != nil {
		fmt.Println("error, cannot create self znode:" + err.Error())
		return
	} 
	
	// addWatch需要的节点
	go listenNode("/gosvr")
}

// stat的数据内容：（在zkCli.sh里面stat或get -s可以看到一部分）
// Czxid 创建节点的事务zxid
// Mzxid 节点最后更新的事务zxid
// Ctime 节点创建时间戳，ms。根节点（"/"）会返回0（即1970-01-01 08:00:00）
// Mtime 节点上次修改时间戳，ms。
// Version 即DataVersion，节点数据版本（每对节点一次set操作会+1）
// Cversion （这个值似乎是每次golang服务器收到一次关于
// Aversion
// EphemeralOwner 如果是临时节点，为节点拥有者的session id。如果不是临时节点，此值为0
// DataLength 节点数据的长度
// NumChildren 节点的子节点数量
// Pzxid 最后一次更新其下子节点（创建删除）的事务ID（自身子节点再往下的节点变更的事务ID不会记在自己的Pzxid身上。子节点发生值变化（如set操作子节点）不会导致自身的Pzxid变化）。自身节点创建时，Pzxid会设为自身创建时用的zxid（即Czxid）
