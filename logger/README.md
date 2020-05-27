## 特性

- 支持格式微自定义
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
