package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

//日志级别
const (
	Debug = 0
	Info  = 1
	Error = 2
)

//logger对象
type logger struct {
	//日志路径
	path string

	//日志文件
	file string

	//打印级别 支持：0-debug/1-info/2-error
	level int

	//日志文件分割大小
	size int64
}

//logger实例的logPrintln方法
func (l *logger) logPrintln(level int, v ...interface{}) {
	//判断是否需要输出日志
	if level < l.level {
		return
	}

	//打开日志文件
	fp, err := l.getFp()
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	//生成日志行首描述字符
	flag := l.getFlag(level)

	//生成调用文件名和行号的简短形式
	short, err := l.getShort(3)
	if err != nil {
		panic(err)
	}

	//重组
	logBuf := []interface{}{short}
	logBuf = append(logBuf, v...)

	//写日志
	logger := log.New(fp, flag, log.Ldate|log.Ltime)
	logger.Println(logBuf...)
}

//logger实例的logPrintf方法
func (l *logger) logPrintf(level int, format string, v ...interface{}) {
	//判断是否需要输出日志
	if level < l.level {
		return
	}

	//打开日志文件
	fp, err := l.getFp()
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	//生成日志行首描述字符
	flag := l.getFlag(level)

	//生成调用文件名和行号的简短形式
	short, err := l.getShort(3)
	if err != nil {
		panic(err)
	}

	//重组
	format = short + " [" + format + "]"

	//写日志
	logger := log.New(fp, flag, log.Ldate|log.Ltime)
	logger.Printf(format, v...)
}

//打开日志文件
//调用者需关闭文件描述符
func (l *logger) getFp() (fp *os.File, err error) {
	//判断日志路径是否存在，不存在则新建
	if len(l.path) > 0 {
		if err = os.MkdirAll(l.path, 0775); err != nil {
			return
		}
	}

	//生成日志文件绝对路径文件名
	f := l.file
	if len(f) <= 0 {
		f = "out.log"
	}
	if len(l.path) > 0 {
		if l.path[len(l.path)-1] == '/' {
			f = l.path + l.file
		} else {
			f = l.path + "/" + l.file
		}
	}

	//日志文件分割
	if err = l.splitFile(f); err != nil {
		return
	}

	//打开文件，不存在则新建
	fp, err = os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	return
}

//日志文件分割
//file为日志文件绝对路径文件名
func (l *logger) splitFile(file string) error {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return err
	}

	//日志文件分割大小，默认50*1024*1024字节
	if l.size <= 0 {
		l.size = 50 * 1024 * 1024
	}

	//日志分割
	if fileInfo.Size() >= l.size {
		t := time.Now()
		date := fmt.Sprintf("%04d%02d%02d", t.Year(), t.Month(), t.Day())
		time := fmt.Sprintf("%02d%02d%02d", t.Hour(), t.Minute(), t.Second())
		tmp := file + "." + date + "." + time
		if err := os.Rename(file, tmp); err != nil {
			return err
		}
	}
	return nil
}

//生成日志行首描述字符
func (l *logger) getFlag(level int) (flag string) {
	switch {
	case level == Debug:
		flag = "[Debug] "
	case level == Info:
		flag = "[Info] "
	case level == Error:
		flag = "[Error] "
	default:
		flag = "[Info] "
	}
	return
}

//生成调用文件名和行号的简短形式
func (l *logger) getShort(skip int) (short string, err error) {
	//获取调用文件名和行号
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		err = fmt.Errorf("runtime.Caller fail")
		return
	}

	//生成简短形式
	short = file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	short = short + ":" + strconv.Itoa(line)
	return
}

//默认logger实例
var l = &logger{}

//设置日志路径
func SetLogPath(path string) {
	l.path = path
}

//设置日志文件
func SetLogFile(file string) {
	l.file = file
}

//设置打印级别
func SetLogLevel(level int) {
	l.level = level
}

//logger实例的Println函数
func Println(level int, v ...interface{}) {
	l.logPrintln(level, v)
}

//logger实例的Printf函数
func Printf(level int, format string, v ...interface{}) {
	l.logPrintf(level, format, v)
}
