package db

import (
	"fmt"
	"strconv"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type MysqlPool struct {
	Inited	bool
	Db		*sql.DB
}

type MysqlConf struct {
	Username	string
	Password	string
	RemoteIp	string
	RemotePort	int
	DbName		string
	OpenConns	int
	IdleConns	int
}

func NewMysqlPool() MysqlPool {
	return MysqlPool {
		Inited: false,
		Db: nil,
	}
}

func (mysql *MysqlPool) InitMysqlPool(conf MysqlConf) {
	if mysql.Inited {
		fmt.Println("InitMysqlPool failed: Mysql Inited");
		return
	}

	var err error
	mysql.Db, err = sql.Open("mysql", conf.Username + ":" + conf.Password + "@tcp(" + conf.RemoteIp + ":" + strconv.Itoa(conf.RemotePort) + ")/" + conf.DbName)
	if err != nil {
		fmt.Println("Init Mysql Error! " + err.Error());
		return
	}
	mysql.Db.SetMaxOpenConns(conf.OpenConns)
	mysql.Db.SetMaxIdleConns(conf.IdleConns)
	mysql.Inited = true
}

func (mysql *MysqlPool) ReleaseMysqlPool() {
	if !mysql.Inited {
		fmt.Println("ReleaseMysqlPool failed: Mysql not inited");
		return
	}

	mysql.Db.Close()
	mysql.Inited = false
}

var db MysqlPool = NewMysqlPool()

func LinkMysql(conf MysqlConf) {

	db.InitMysqlPool(conf)
	defer db.ReleaseMysqlPool()

	if db.Inited {
		stmt, _ := db.Db.Prepare(`insert into test_db.test_table (value) values (?)`)
		defer stmt.Close()
	
		rows, err := stmt.Query("hello")
		defer rows.Close()
	
		if err != nil {
			fmt.Println("insert error: " + err.Error())
			return
		}
	}
}