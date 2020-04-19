package main

import (
	"encoding/json"
	"fmt"
	"logger"
)

type I1001 struct {
	TransCode string
	Test      string
}

type O1001 struct {
	Ret  int    `json:"ret"`
	Msg  string `json:"msg"`
	Test string `json:"test"`
}

var i1001 I1001
var o1001 O1001

func p1001(inBuf []byte) (outBuf []byte) {
	//入参检查
	if nil == inBuf || len(inBuf) <= 0 {
		logger.Println(logger.Error, "非法入参")
		return
	}

	//解析请求
	if err := json.Unmarshal(inBuf, &i1001); err != nil {
		logger.Println(logger.Error, err)
		return
	}
	logger.Println(logger.Info, "解析请求成功")

	//业务处理
	logger.Println(logger.Info, "业务处理开始")
	logger.Println(logger.Info, i1001.TransCode, i1001.Test)
	logger.Println(logger.Info, "业务处理结束")

	//构造应答
	o1001.Ret = 0
	o1001.Msg = "succ"
	o1001.Test = fmt.Sprintf("hello %s %s", i1001.TransCode, i1001.Test)
	outBuf, err := json.Marshal(&o1001)
	if err != nil {
		logger.Println(logger.Error, err)
		return
	}
	logger.Println(logger.Info, "构造应答成功")
	return
}
