package comm

import (
	"encoding/json"
	"logger"
)

//router对象
type router struct {
	route map[string]func([]byte) []byte
}

//router实例handle方法
func (r *router) handle(inBuf []byte) (outBuf []byte) {
	//入参检查
	if nil == inBuf || len(inBuf) <= 0 {
		panic("非法入参")
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
	logger.Printf(logger.Info, "交易码%s", t.TransCode)

	//调用业务处理函数
	h, ok := r.route[t.TransCode]
	if !ok {
		outBuf = []byte("非法交易")
	}
	outBuf = h(inBuf)
	return
}
