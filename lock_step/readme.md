# 帧同步

项目并没有在任何地方使用帧同步，因为并没有实时性要求特别高的场景，血量蓝量什么的数据都是后端算好之核每1秒给镜头内的客户端同步一次就行。

（其实某个特殊场合这种纯状态同步的方式导致一部分玩家特别卡，但最后也没引入帧同步做这部分优化，可能认为没必要或者怕如果提出搞这个的话客户端当场罢工罢，又或者问题不是出在同步方式上）
