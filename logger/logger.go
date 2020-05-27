package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
)

/*
	保留日志天数，定时扫描 dir
*/
type Logger struct {
	level  int32
	cfg    Config
	format *Format

	// 写入
	write *Write
}

const (
	DebugLevel int = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel

	// 短方法名长度
	_shortFuncname = 2
)

// Config 日志配置信息
type Config struct {
	// 调用深度
	Calldpeth int
	// 日志等级，>= 设置的等级才会写入
	Level int
	// 格式化，配置
	Format *FormatConfig
	Write  *WriteConfig
}

// DefaultConfig 默认配置
func DefaultConfig() Config {
	return Config{
		Calldpeth: 2,
		Level:     DebugLevel,
		Format:    DefaultFormatConfig(),
		Write:     DefaultWriteConfig(),
	}
}

// NewLogger
func NewLogger(cfg Config) *Logger {
	lg := Logger{
		level:  int32(cfg.Level),
		cfg:    cfg,
		format: NewFormat(cfg.Format),
	}
	return &lg
}

// Debug
func (lg *Logger) Debug(format string, args ...interface{}) {
	if !lg.isWrite(DebugLevel) {
		return
	}
	lg.Append(lg.String(lg.cfg.Calldpeth, "DEBUG", fmt.Sprintf(format, args...)))
}

// Info
func (lg *Logger) Info(format string, args ...interface{}) {
	if !lg.isWrite(InfoLevel) {
		return
	}
	lg.Append(lg.String(lg.cfg.Calldpeth, "INFO", fmt.Sprintf(format, args...)))
}

// Warn
func (lg *Logger) Warn(format string, args ...interface{}) {
	if !lg.isWrite(WarnLevel) {
		return
	}
	lg.Append(lg.String(lg.cfg.Calldpeth, "WARN", fmt.Sprintf(format, args...)))
}

// Error
func (lg *Logger) Error(format string, args ...interface{}) {
	if !lg.isWrite(ErrorLevel) {
		return
	}
	lg.Append(lg.String(lg.cfg.Calldpeth, "ERROR", fmt.Sprintf(format, args...)))
}

// Fatal 等级 退出程序
func (lg *Logger) Fatal(format string, args ...interface{}) {
	if !lg.isWrite(FatalLevel) {
		return
	}
	lg.Append(lg.String(lg.cfg.Calldpeth, "FATAL", fmt.Sprintf(format, args...)))
	os.Exit(1)
}

// String 生成待写入文件的数据
func (lg *Logger) String(calldpeth int, level, message string) []byte {

	if calldpeth > 0 {
		_, file, line, ok := runtime.Caller(calldpeth)
		if !ok {
			file = "???"
			line = 0
		}

		temp := strings.Split(file, "/")
		if len(temp) >= _shortFuncname {
			temp = temp[len(temp)-_shortFuncname:]
		}
		file = strings.Join(temp, "/")

		return lg.format.GenMessage(level, file+":"+strconv.Itoa(line)+"\t"+message)
	}

	return lg.format.GenMessage(level, message)
}

// Append 写入文件
func (lg *Logger) Append(message []byte) (n int, err error) {
	if lg.write == nil {

		lg.write = NewWrite(lg.cfg.Write)
	}
	return lg.write.Write(message)
}

// SetLevel 设置错误等级
func (lg *Logger) SetLevel(level int) {
	atomic.StoreInt32(&lg.level, int32(level))
}

// Stack 堆栈信息
func (lg *Logger) Stack(err error) string {
	return lg.format.Stack(err.Error())
}

// Close
func (lg *Logger) Close() error {
	if lg.write != nil {
		return lg.write.Close()
	}
	return nil
}

func (lg *Logger) isWrite(level int) bool {
	return int32(level) >= atomic.LoadInt32(&lg.level)
}

var stdout io.Writer = os.Stderr

// Output 系统输出
func (lg *Logger) Output(calldepth int, s string) (n int, err error) {
	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0
	}
	return stdout.Write(lg.format.GenMessage("output", file+":"+strconv.Itoa(line)+" "+s))
}
