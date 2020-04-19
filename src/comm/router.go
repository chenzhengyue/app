package comm

import (
	"encoding/json"
	"logger"
)

//handleFunc类型
type handleFunc func([]byte) []byte

//router对象
type router struct {
	route map[string]handleFunc
}

//router实例handle方法
func (r *router) handle(inBuf []byte) (outBuf []byte) {
	//入参检查
	if nil == inBuf || len(inBuf) <= 0 {
		logger.Println(logger.Error, "非法入参")
		return
	}

	//交易码数据结构
	type transcode struct {
		TransCode string
	}

	//解析交易码
	var t transcode
	if err := json.Unmarshal(inBuf, &t); err != nil {
		outBuf = []byte("解析交易码错")
		logger.Println(logger.Error, err)
		return
	}
	logger.Printf(logger.Info, "交易码：%s\n", t.TransCode)

	//调用业务处理函数
	if h, ok := r.route[t.TransCode]; !ok {
		outBuf = []byte("非法交易")
	} else {
		outBuf = h(inBuf)
	}
	return
}
