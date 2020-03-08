package mysql

import (
	"database/sql"
	"fmt"
	"github.com/reptile/config"
	_ "github.com/reptile/dependency_pack/go-sql-driver"
	"os"
)

var db *sql.DB

func init() {
	dbInfo := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", config.MysqlDbInfo.UserName, config.MysqlDbInfo.Password, config.MysqlDbInfo.Addr, config.MysqlDbInfo.DataBaseName)
	db, _ = sql.Open("mysql", dbInfo)
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql, err:" + err.Error())
		os.Exit(1)
	}
}

func DBCon() *sql.DB {
	return db
}
