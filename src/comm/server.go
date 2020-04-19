package comm

import (
	"fmt"
	"logger"
	"net"
	"strconv"
	"time"
)

//handler接口
type handler interface {
	handle([]byte) []byte
}

//默认handler
var defaultHandler = &router{make(map[string]handleFunc)}

//默认handler注册函数
func HandleFunc(k string, f func([]byte) []byte) {
	//入参检查
	if "" == k || nil == f {
		logger.Println(logger.Error, "非法入参")
		return
	}

	//注册函数
	if _, ok := defaultHandler.route[k]; ok {
		logger.Printf(logger.Error, "重复注册")
		return
	} else {
		defaultHandler.route[k] = f
	}
}

//Server对象
type Server struct {
	//地址：":8080"
	Addr string

	//handler接口
	Handler handler

	//连接超时时间 默认30s
	ConnTimeout time.Duration

	//读报文头超时时间 默认5s
	ReadHeaderTimeout time.Duration

	//读报文体超时时间 默认不设置，和读报文头共用5s
	ReadBodyTimeout time.Duration

	//写应答超时时间 默认5s
	WriteTimeout time.Duration
}

//Server实例启动函数
func Run(addr ...string) {
	//入参检查
	if len(addr) > 1 {
		logger.Println(logger.Error, "调用方法：Run(addr)，addr选送")
		return
	}

	//实例化Server
	s := &Server{}

	//赋值Addr
	if 1 == len(addr) {
		s.Addr = addr[0]
	}

	//启动Server实例
	s.Run()
}

//Server实例启动方法
func (s *Server) Run() {
	for {
		s.serve()
		logger.Println(logger.Info, "Server restart")
	}
}

//Server实例启动方法
func (s *Server) serve() {
	//异常恢复，服务不停机
	defer func() {
		if err := recover(); err != nil {
			logger.Println(logger.Error, err)
		}
	}()

	//创建监听
	if "" == s.Addr {
		s.Addr = ":8080"
	}
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		logger.Println(logger.Error, err)
		return
	}
	defer l.Close()
	logger.Println(logger.Info, "listen running on", s.Addr)

	//accept失败后，重新accept的等待时间
	tempDelay := 0 * time.Second
	//accept失败时，重新accept的最长等待时间
	maxDelay := 1 * time.Second

	//处理连接
	for {
		con, err := l.Accept()
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Temporary() {
				if 0 == tempDelay {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if tempDelay > maxDelay {
					tempDelay = maxDelay
				}
				logger.Println(logger.Error, err)
				logger.Printf(logger.Error, "wait [%v], then accept again\n", tempDelay)
				time.Sleep(tempDelay)
				continue
			} else {
				logger.Println(logger.Error, err)
				return
			}
		}

		//实例化conn
		c := &conn{server: s, rwc: con}
		//启动协程处理conn
		//暂不控制协程最大数量
		//绝大部分场景可正常服务，极端情况需启动异常多的协程且明显影响系统性能时，再做优化
		go c.serve()
	}
}

//conn对象
type conn struct {
	//继承Server对象
	server *Server

	//net.Conn对象
	rwc net.Conn

	//请求buf
	inBuf []byte

	//应答buf
	outBuf []byte
}

//conn实例serve方法
func (c *conn) serve() {
	defer func() {
		//关闭连接
		c.rwc.Close()

		//异常恢复，不影响其它协程
		if err := recover(); err != nil {
			logger.Println(logger.Error, err)
		}
	}()

	//设置连接超时
	if 0 == c.server.ConnTimeout {
		c.server.ConnTimeout = 30 * time.Second
	}
	t := time.Now().Add(c.server.ConnTimeout)
	if err := c.rwc.SetDeadline(t); err != nil {
		logger.Println(logger.Error, err)
		return
	}

	//通讯开始时间
	start := time.Now()
	logger.Println(logger.Info, "接收客户端请求开始")

	//获取客户端信息
	logger.Println(logger.Info, "客户端信息：", c.rwc.RemoteAddr())

	//读取客户端请求
	if err := c.readConn(); err != nil {
		logger.Println(logger.Error, err)
		return
	}

	//处理客户端请求
	if nil == c.server.Handler {
		c.server.Handler = defaultHandler
	}
	c.outBuf = c.server.Handler.handle(c.inBuf)

	//应答客户端请求
	if err := c.writeConn(); err != nil {
		logger.Println(logger.Error, err)
		return
	}

	//通讯结束时间
	end := time.Now()
	logger.Println(logger.Info, "应答客户端请求结束")

	//通讯耗时
	elapsed := end.Sub(start)
	logger.Println(logger.Info, "conn use time:", elapsed)
}

//conn实例读方法
func (c *conn) readConn() error {
	//设置报文头读超时
	if 0 == c.server.ReadHeaderTimeout {
		c.server.ReadHeaderTimeout = 5 * time.Second
	}
	t := time.Now().Add(c.server.ReadHeaderTimeout)
	if err := c.rwc.SetReadDeadline(t); err != nil {
		return err
	}

	//通讯规范：报文头 + 报文体
	//通讯规范：报文头：4字节，表示报文体的长度；不足4位，左补空格/右补空格/左补0
	//通讯规范：报文体：最大长度支持 1024*9 字节
	//通讯规范：报文体格式：通讯中不区分报文体格式，业务逻辑中处理
	//通讯规范：报文编码：utf-8
	//接收报文头
	const headLen = 4
	recvLen, try := 0, 0
	inBuf := make([]byte, headLen)
	for recvLen < headLen {
		//已验证不会读超过4个字节
		n, err := c.rwc.Read(inBuf[recvLen:])
		if err != nil {
			//因设置了读超时，读超时导致net.Error错误，且为临时性错误，因此设置try
			if e, ok := err.(net.Error); ok && e.Temporary() && try < 5 {
				logger.Println(logger.Error, err)
				logger.Println(logger.Error, "临时性错误，重试")
				try++
				//此处不使用continue，防止报错但读到了内容
			} else {
				logger.Println(logger.Error, "重试次数超限，退出")
				return err
			}
		}
		recvLen += n
	}
	bodyLen, err := strconv.Atoi(string(inBuf))
	if err != nil {
		return err
	}
	logger.Println(logger.Info, "recv head:", bodyLen)

	//校验报文体长度
	const bodyMaxLen = 1024 * 9
	if bodyLen > bodyMaxLen || bodyLen <= 0 {
		//设置写超时
		if 0 == c.server.WriteTimeout {
			c.server.WriteTimeout = 5 * time.Second
		}
		t := time.Now().Add(c.server.WriteTimeout)
		if err := c.rwc.SetWriteDeadline(t); err != nil {
			return err
		}

		if bodyLen <= 0 {
			//准备应答数据
			outBuf := []byte("body len in head is zero")

			//发送应答数据
			_, err := c.rwc.Write(outBuf)
			if err != nil {
				return err
			}
			logger.Println(logger.Info, "outBuf:", string(outBuf))

			//返回
			return fmt.Errorf("body len in head err, bodyLen[%d]", bodyLen)
		} else {
			//准备应答数据
			outBuf := []byte("body len in head, too long")

			//发送应答数据
			_, err := c.rwc.Write(outBuf)
			if err != nil {
				return err
			}
			logger.Println(logger.Info, "outBuf:", string(outBuf))

			//返回
			return fmt.Errorf("body len in head, too long, bodyLen[%d] bodyMaxLen[%d]", bodyLen, bodyMaxLen)
		}
	}

	//设置报文体读超时
	if 0 != c.server.ReadBodyTimeout {
		t = time.Now().Add(c.server.ReadBodyTimeout)
		if err := c.rwc.SetReadDeadline(t); err != nil {
			return err
		}
	}

	//接收报文体
	recvLen, try = 0, 0
	//根据报文头中说明的报文体长度，动态分配buf长度，不会造成资源浪费
	c.inBuf = make([]byte, bodyLen)
	for recvLen < bodyLen {
		n, err := c.rwc.Read(c.inBuf[recvLen:])
		if err != nil {
			//因设置了读超时，读超时导致net.Error错误，且为临时性错误，因此设置try
			if e, ok := err.(net.Error); ok && e.Temporary() && try < 5 {
				logger.Println(logger.Error, err)
				logger.Println(logger.Error, "临时性错误，重试")
				try++
				//此处不使用continue，防止报错但读到了内容
			} else {
				logger.Println(logger.Error, "重试次数超限，退出")
				return err
			}
		}
		recvLen += n
	}
	logger.Println(logger.Info, "recv body:", string(c.inBuf))
	return nil
}

//conn实例写方法
func (c *conn) writeConn() error {
	//判断是否存在应答数据
	if len(c.outBuf) <= 0 {
		return fmt.Errorf("outBuf no data, outBuf[%s]", c.outBuf)
	}

	//设置写超时
	if 0 == c.server.WriteTimeout {
		c.server.WriteTimeout = 5 * time.Second
	}
	t := time.Now().Add(c.server.WriteTimeout)
	if err := c.rwc.SetWriteDeadline(t); err != nil {
		return err
	}

	//写应答数据
	sendLen, err := c.rwc.Write(c.outBuf)
	if err != nil {
		return err
	} else if sendLen < len(c.outBuf) {
		return fmt.Errorf("data writed not enough, sendLen[%d] bufLen[%d]", sendLen, len(c.outBuf))
	}
	logger.Println(logger.Info, "sendLen:", sendLen)
	logger.Println(logger.Info, "outBuf:", string(c.outBuf))
	return nil
}
