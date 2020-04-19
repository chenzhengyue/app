package main

import (
	"github.com/chenzhengyue/comm"
	"github.com/chenzhengyue/logger"
	"os"
)

func init() {
	path := os.Getenv("HOME") + "/app/log"
	logger.SetLogPath(path)
	logger.SetLogFile("main.log")
	logger.SetLogLevel(logger.Info)
}

func main() {
	//路由注册
	comm.HandleFunc("1001", p1001)

	//启动服务
	comm.Run()
}
