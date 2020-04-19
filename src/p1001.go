package main

import (
	"encoding/json"
	"fmt"
	"github.com/chenzhengyue/logger"
)

type I1001 struct {
	TxnCode string
	Test    string
}

type O1001 struct {
	Ret  int    `json:"ret"`
	Msg  string `json:"msg"`
	Test string `json:"test"`
}

var pi1001 I1001
var po1001 O1001

func p1001(inBuf []byte) (outBuf []byte) {
	//入参检查
	if nil == inBuf || len(inBuf) <= 0 {
		logger.Println(logger.Error, "非法入参")
		return
	}

	//解析请求
	if err := json.Unmarshal(inBuf, &pi1001); err != nil {
		logger.Println(logger.Error, err)
		return
	}
	logger.Println(logger.Info, "解析请求成功")

	//业务处理
	logger.Println(logger.Info, "业务处理开始")
	logger.Println(logger.Info, pi1001.TxnCode, pi1001.Test)
	logger.Println(logger.Info, "业务处理结束")

	//构造应答
	po1001.Ret = 0
	po1001.Msg = "succ"
	po1001.Test = fmt.Sprintf("hello %s %s", pi1001.TxnCode, pi1001.Test)
	outBuf, err := json.Marshal(&po1001)
	if err != nil {
		logger.Println(logger.Error, err)
		return
	}
	logger.Println(logger.Info, "构造应答成功")
	return
}
