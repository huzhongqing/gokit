# simplelog

[![Build Status](https://travis-ci.org/huzhongqing/simplelog.svg?branch=master)](https://travis-ci.org/huzhongqing/simplelog)
[![Go Report Card](https://goreportcard.com/badge/github.com/huzhongqing/simplelog)](https://goreportcard.com/report/github.com/huzhongqing/simplelog)
[![codecov](https://codecov.io/gh/huzhongqing/simplelog/branch/master/graph/badge.svg)](https://codecov.io/gh/huzhongqing/simplelog)

## 特性

- 支持格式微自定义
- 不同错误级别写入不同文件夹
- 单文件容量限制,切割文件
- 文件保留时间设置，自动删除过期文件
- 捕获堆栈信息
 
## Usage
```
go get -u github.com/huzhongqing/gokit/logger
```

``` go
func TestLogger(t *testing.T) {
	lg := NewLogger(DefaultConfig())

	lg.Debug("DEBUG")
	lg.Info("INFO")
	lg.Warn("WARN")
	lg.Error("ERROR")
	lg.Error("%s_%s", "stack", lg.Stack(fmt.Errorf("这是stack")))
	lg.Fatal("FATAL")

	_ = lg.Close()
}
```
