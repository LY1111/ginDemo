package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"tag_data_sync/config"
	"time"
)

var Log  *zap.SugaredLogger

func InitLogger(conf *config.LogConfig) error {
	isDebug := true
	if conf.IsDebug != 1 {
		isDebug = false
	}
	err := initLogger(conf.LogFile, conf.LogLevel, isDebug, conf.MaxSize, conf.MaxBackups, conf.MaxAge, conf.Compress)
	if err != nil {
		return err
	}
	log.SetFlags(log.Lmicroseconds | log.Lshortfile | log.LstdFlags)
	return nil
}

func ZnTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func initLogger(logFile string, logLevel string, isDebug bool, maxSize, maxBackups, maxAge int, compress bool) error {
	hook := lumberjack.Logger{
		Filename:   logFile,    // ⽇志⽂件路径
		MaxSize:    maxSize,    // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: maxBackups, // 最多保留3个备份
		MaxAge:     maxAge,     // 保留旧文件的最大天数
		Compress:   compress,   // 是否压缩 disabled by default
		LocalTime:  true,
	}
	fileWriter := zapcore.AddSync(&hook)

	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)

	// for human operators.
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = ZnTimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	//初始化所有core
	var allCore []zapcore.Core

	if isDebug {
		allCore = append(allCore, zapcore.NewCore(consoleEncoder, consoleDebugging, level))
	}

	allCore = append(allCore, zapcore.NewCore(consoleEncoder, fileWriter, level))

	core := zapcore.NewTee(allCore...)

	// From a zapcore.Core, it's easy to construct a Logger.
	zlog := zap.New(core).WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))

	Log = zlog.Sugar()
	return nil
}

// Debug fmt.Sprintf to log a templated message.
func Debug(args ...interface{}) {
	argsMerge := AddHostName(args...)
	Log.Debug(argsMerge...)
}

// Info uses fmt.Sprintf to log a templated message.
func Info(args ...interface{}) {
	argsMerge := AddHostName(args...)
	Log.Info(argsMerge...)
}

// Warn uses fmt.Sprintf to log a templated message.
func Warn(args ...interface{}) {
	argsMerge := AddHostName(args...)
	SendMonitor2DingDing(args)
	Log.Warn(argsMerge...)
}

// Error uses fmt.Sprintf to log a templated message.
func Error(args ...interface{}) {
	argsMerge := AddHostName(args...)
	SendMonitor2DingDing(args)
	Log.Error(argsMerge...)
}

// Fatal uses fmt.Sprintf to log a templated message.
func Fatal(args ...interface{}) {
	SendMonitor2DingDing(args)
	Log.Fatal(args...)
}

// Debugf fmt.Sprintf to log a templated message.
func Debugf(format string, args ...interface{}) {
	Log.Debugf(GetHostNamePrefix()+format, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(format string, args ...interface{}) {
	Log.Infof(GetHostNamePrefix()+format, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	go SendMonitor2DingDing(str)
	Log.Warnf(GetHostNamePrefix()+format, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	go SendMonitor2DingDing(str)
	Log.Errorf(GetHostNamePrefix()+"   "+format, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message.
func Fatalf(format string, args ...interface{}) {
	Log.Fatalf(format, args...)
}

/**
* @Author liyan
* @Description  钉钉报警（只有生产环境才会发送钉钉信息）
* @Date   2020/11/12 2:09 下午
* @Param
* @return
**/
func SendMonitor2DingDing(args ...interface{}) {
	// 正式环境才会报警
	//if config.Config.Service.Env != "pro" {
		return
	//}

	slice := make([]string, len(args))

	for i, v := range args {
		slice[i] = fmt.Sprint(v)
	}

	s := strings.Join(slice, ",")

	host, _ := os.Hostname()

	 b := json.RawMessage(`
		{"msgType": "text","text": {"content": "bussiness: [` + host + `-` + "" + `] 数据桥接业务报警\n ` + s + `"}}
	 `)

	url := ""
	//http.Post(url, "application/json", strings.NewReader(string(b))) //忽略dingding错误

	//PostJson(url, []byte(b))

	fmt.Println(b,url)
}

func AddHostName(args ...interface{}) []interface{} {
	argsMerge := make([]interface{}, 0, len(args)+1)
	argsMerge = append(argsMerge, GetHostNamePrefix())
	argsMerge = append(argsMerge, args...)
	return argsMerge
}

// HttpPost post请求
func PostJson(url string, params []byte) ([]byte, error) {
	body := bytes.NewBuffer(params)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("httpcode error:" + fmt.Sprint(resp.StatusCode))
	}

	return respData, nil
}

//获取本机hostname
var HostName = "CCServer"

func InitHostName() error {
	HostName, _ = os.Hostname()
	return nil
}

func GetHostNamePrefix() string {
	InitHostName()
	return HostName + "	"
}
