package user

import (
	"go_svr/log"
	"go_svr/redis"
)

type User struct {
	UserId uint64
	OpenId string
	Name   string
	Dirty  bool // 这个玩家是否产生了数据变更
	// 别的
}

type UserManager struct {
	Users     map[uint64]*User
	Open2User map[string]*User // 一般而言 一个openid只能创建1个角色，所以可以这么索引
	//Name2User map[string]*User	// 是否允许名字重复这个，有些游戏是允许的，有些不能，因此这个索引只在不允许重复的时候可以创建（可能仅仅是方便运营问名字定位玩家而设定的，现在uid一般都恒定镶在游戏界面的某个角落了）
	DirtyUser map[uint64]*User
}

var m = &UserManager{
	Users:     make(map[uint64]*User),
	Open2User: make(map[string]*User),
	DirtyUser: make(map[uint64]*User),
}

func GetMgr() *UserManager {
	return m
}

func (um *UserManager) Init() {

}

func (um *UserManager) GetUser(userId uint64) *User {
	return um.Users[userId]
}
func (um *UserManager) GetUserByOpen(openId string) *User {
	return um.Open2User[openId]
}

func (um *UserManager) CreateUser(openId string) *User {
	uid := redis.GetRedisCli().Incr("MaxUserId") // todo 也有可能使用snowflake生成唯一ID
	u := &User{
		UserId: uint64(uid),
		OpenId: openId,
		Name:   "",
		Dirty:  true,
	}
	um.Users[uint64(uid)] = u
	um.Open2User[openId] = u
	um.DirtyUser[uint64(uid)] = u
	log.Info("User %d created, open id = %s", u.UserId, openId)
	return u
}

func (um *UserManager) SetDirty(userId uint64) {
	u, ok := um.Users[userId]
	if !ok {
		return
	}
	u.Dirty = true
	um.DirtyUser[userId] = u
}
