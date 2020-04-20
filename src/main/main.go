package main

import (
	"comm"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"logger"
	"os"
)

func init() {
	path := os.Getenv("HOME") + "/app/log"
	logger.SetLogPath(path)
	logger.SetLogLevel(logger.Info)
	logger.SetLogFile("main.log")
}

//数据库实例
var db *sql.DB

func main() {
	//异常处理
	defer func() {
		if err := recover(); err != nil {
			logger.Println(logger.Error, err)
		}
	}()

	//连接数据库
	var err error
	db, err = sql.Open("mysql", "go:%TGB6yhn@tcp(192.168.234.128:3306)/go")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//路由注册
	if err = setRouter(); err != nil {
		panic(err)
	}

	//启动服务
	comm.Run()
}

//路由注册
func setRouter() error {
	//查询路由表
	rows, err := db.Query("select trans_code, handler from tbl_routers where flag = ?", 1)
	if err != nil {
		return err
	}
	defer rows.Close()

	//读取路由表
	var transCode string
	var handler string
	for rows.Next() {
		if err := rows.Scan(&transCode, &handler); err != nil {
			return err
		}
		logger.Println(logger.Info, transCode, handler)

		//注册路由
		comm.HandlerFunc(transCode, p1001)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}
