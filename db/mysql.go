package db

import (
	"fmt"
	"log"
	"strconv"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type MysqlPool struct {
	Inited bool
	Db     *sql.DB
}

type MysqlConf struct {
	Username   string
	Password   string
	RemoteIp   string
	RemotePort int
	DbName     string
}

func NewMysqlPool() *MysqlPool {
	return &MysqlPool{
		Inited: false,
		Db:     nil,
	}
}

func (mysql *MysqlPool) InitMysqlPool(conf MysqlConf) {
	if mysql.Inited {
		fmt.Println("InitMysqlPool failed: Mysql Inited")
		return
	}

	var err error
	mysql.Db, err = sql.Open("mysql", conf.Username+":"+conf.Password+"@tcp("+conf.RemoteIp+":"+strconv.Itoa(conf.RemotePort)+")/"+conf.DbName)
	if err != nil {
		panic("Init Mysql Error! " + err.Error())
	}
	mysql.Db.SetMaxOpenConns(1)
	mysql.Db.SetMaxIdleConns(1)
	mysql.Inited = true
	log.Printf("init mysql pool success")
}

func (mysql *MysqlPool) OnClose() {
	if !mysql.Inited {
		fmt.Println("ReleaseMysqlPool failed: Mysql not inited")
		return
	}

	mysql.Inited = false
	err := mysql.Db.Close()
	if err != nil {
		panic("close mysql error: " + err.Error())
	}
	log.Printf("release mysql pool success")
}

var db = NewMysqlPool()

func GetDbPool() *MysqlPool {
	return db
}

func (mysql *MysqlPool) QueryMySql(sql string, args ...any) error {
	if !mysql.Inited {
		panic("QueryMySql failed, mysql is not inited")
	}

	rows, err := db.Db.Query(sql, args...)
	if err != nil {
		log.Printf("mysql query error: %s", err.Error())
		return err
	}
	defer rows.Close()
	return nil
}
