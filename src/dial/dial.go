package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	//调用参数检查
	if len(os.Args) < 3 {
		fmt.Println("调用方式：go run dial.go clientNum txnNum")
		return
	}

	//参数初始化
	clientNum, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("strconv.Atoi err", err)
		return
	}
	txnNum, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("strconv.Atoi err", err)
		return
	}

	//开始时间
	start := time.Now()

	quit := make(chan bool)
	for i := 0; i < clientNum; i++ {
		go client(txnNum, quit)
	}
	for i := 0; i < clientNum; i++ {
		<-quit
	}

	//结束时间
	end := time.Now()

	//耗时
	elapsed := end.Sub(start)
	fmt.Println("use time:", elapsed)
}

func client(txnNum int, quit chan bool) {
	defer func() {
		//异常恢复
		err := recover()
		if err != nil {
			fmt.Println("client recover:", err)
		}
		//通知主协程
		quit <- true
	}()

	//处理交易
	for i := 0; i < txnNum; i++ {
		handleDeal()
	}
}

func handleDeal() {
	//异常处理
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("handleDeal recover:", err)
		}
	}()

	//连接服务器
	conn, err := net.DialTimeout("tcp", "192.168.234.128:8080", 5*time.Second)
	if err != nil {
		fmt.Println("net.Dial err", err)
		return
	}
	defer conn.Close()

	//time.Sleep(20 * time.Second)
	//发送数据
	inBuf := "{\"TransCode\":\"1001\",\"test\":\"test\"}"
	inLen := fmt.Sprintf("%04d", len(inBuf))
	data := inLen + inBuf
	_, err = conn.Write([]byte(data))
	if err != nil {
		fmt.Println("conn.Write err", err)
		conn.Close()
		return
	}

	//接收数据
	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println("conn.Write err", err)
		conn.Close()
		return
	}
	fmt.Println(string(buf))
}
