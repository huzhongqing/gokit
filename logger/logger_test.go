package logger

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
	"time"
)

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

func TestLogger_Size(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Write.Compress = true
	cfg.Write.MaxSize = 1024 * 1024
	lg := NewLogger(cfg)

	count := 100
	for count > 0 {
		count--
		lg.Info("TestSimpleLogger_Size")
		lg.Debug("TestSimpleLogger_Size")
		lg.Warn("TestSimpleLogger_Size")
		lg.Error("TestSimpleLogger_Size")
	}

	lg.Close()
}

func TestLogger_Keep(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Write.MaxAge = 10 * time.Second
	lg := NewLogger(cfg)

	lg.Info("info")

	lg.Close()
}

func TestLogger_SetLevel(t *testing.T) {
	lg := NewLogger(DefaultConfig())

	lg.SetLevel(ErrorLevel)

	lg.Debug("TestSimpleLogger_SetLevel")

	lg.Close()

	content, _ := ioutil.ReadFile(lg.cfg.Write.Filename)
	if strings.Index(string(content), "TestSimpleLogger_SetLevel") != -1 {
		t.Fail()
	}

}

func TestFormat_Stack(t *testing.T) {
	f := NewFormat(DefaultFormatConfig())
	if strings.Index(f.Stack("-"), "Traceback:") != -1 {
		return
	}
	t.Fatal()
}

func TestFormat_GenMessage(t *testing.T) {
	f := NewFormat(DefaultFormatConfig())
	t.Log(f.GenMessage("DEBUG", "TestFormat_GenMessage"))
}

var (
	dir    = "./testdata"
	maxAge = time.Hour
)

func TestWrite_isRoll(t *testing.T) {

	w := NewWrite(&WriteConfig{
		Filename: dir + "/TestWrite_isRoll.log",
		MaxSize:  50,
		MaxAge:   maxAge,
		Compress: true,
	})

	count := 30
	for count > 0 {
		count--
		w.Write([]byte("size test" + strconv.Itoa(count)))
	}

	time.Sleep(500 * time.Millisecond)

	err := w.Close()
	if err != nil {
		t.Log(err)
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
		return
	}

	for _, v := range files {
		if strings.Index(v.Name(), "TestWrite_isRoll.2") != -1 {
			t.Log("backup", v.Name())
			return
		}
	}
	t.Fatal("backup error")
}
