package cache

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

/*
	内存实现cache
Notice:
	1. 定期过期，清理申请的内存。
	2. 仅适用一些小的缓存，不适合存储太大的数据量
*/

var (
	ErrKeysOverCapacity = errors.New("keys over capacity")
	ErrFilenameEmpty    = errors.New("save filename empty")
)

type Options struct {
	// 缓存容量， 默认 -1  不限制
	Capacity int
	// 保存文件位置, 默认 不设置，不能 Save
	Filename string
	// 自动清除, 当ttl key存在不多时，可以关闭，默认开启
	AutoClean bool
}

func NewDefaultOptions() Options {
	return Options{
		Capacity: -1,
		Filename: "",
	}
}

type MemCache struct {
	rwMutex sync.RWMutex
	store   map[string]interface{}

	// 缓存容量， -1 - 不限制
	capacity int
	// 保存文件位置, 不设置，默认当前执行路径
	filename string
	// 自动清除
	autoClean bool
}

func NewMemCache(opts ...Options) *MemCache {
	opt := NewDefaultOptions()
	if len(opts) > 0 {
		opt = opts[0]
	}
	mem := &MemCache{
		store:     make(map[string]interface{}),
		capacity:  opt.Capacity,
		filename:  opt.Filename,
		autoClean: opt.AutoClean,
	}

	if mem.autoClean {
		go mem.autoExpireClean(10 * time.Minute)
	}

	if err := mem.load(); err != nil {
		log.Printf("WARNING: load file cache error %s \n", err.Error())
	}

	return mem
}

type ttlValue struct {
	Value       interface{} `json:"v"`
	ExpiredTime time.Time   `json:"e"`
}

func (val *ttlValue) SetExpiredTime(t time.Duration) {
	if t <= -1 {
		val.ExpiredTime = time.Now().AddDate(100, 0, 0)
		return
	}
	val.ExpiredTime = time.Now().Add(t)
}

func (val *ttlValue) TTL() time.Duration {
	expire := val.ExpiredTime.Sub(time.Now())
	// 存在这种可能
	if expire < 0 {
		expire = 0
	}
	return expire
}

func (val *ttlValue) Expired() bool {
	return val.ExpiredTime.Before(time.Now())
}

func (mem *MemCache) Get(key string) *Cmd {
	mem.rwMutex.RLock()
	defer mem.rwMutex.RUnlock()

	v, ok := mem.store[key]
	if ok {
		val := v.(ttlValue)
		// 如果过期了，就删除了
		if val.Expired() {
			mem.rwMutex.Lock()
			delete(mem.store, key)
			mem.rwMutex.Unlock()
			return &Cmd{baseCmd: baseCmd{exists: false, ttl: val.TTL()}, value: nil}
		}

		return &Cmd{baseCmd: baseCmd{exists: ok, ttl: val.TTL()}, value: val.Value}
	}
	return &Cmd{}
}

func (mem *MemCache) Set(key string, value interface{}, ttl time.Duration) *StatusCmd {
	mem.rwMutex.Lock()
	defer mem.rwMutex.Unlock()
	if err := mem.isOverCapacity(); err != nil {
		return &StatusCmd{baseCmd: baseCmd{err: err}}
	}
	val := ttlValue{
		Value: value,
	}
	val.SetExpiredTime(ttl)
	mem.store[key] = val
	return &StatusCmd{baseCmd: baseCmd{exists: true, ttl: val.TTL()}, value: StatusOK}
}

func (mem *MemCache) Delete(key string) *StatusCmd {
	mem.rwMutex.Lock()
	defer mem.rwMutex.Unlock()

	delete(mem.store, key)
	return &StatusCmd{value: StatusOK}
}

func (mem *MemCache) Keys(prefix string) *SliceStringCmd {
	mem.rwMutex.RLock()
	defer mem.rwMutex.RUnlock()
	keys := make([]string, 0)
	for key := range mem.store {
		if strings.HasPrefix(key, prefix) {
			keys = append(keys, key)
		}
	}
	return &SliceStringCmd{value: keys}
}

// FlushAll 清空所有数据
func (mem *MemCache) FlushAll() *StatusCmd {
	mem.rwMutex.Lock()
	defer mem.rwMutex.Unlock()
	mem.store = make(map[string]interface{})
	return &StatusCmd{value: StatusOK}
}

// Close 开启写入磁盘，则写入文件
func (mem *MemCache) Close() error {
	return nil
}

func (mem *MemCache) isOverCapacity() error {
	if mem.Capacity == -1 {
		return nil
	}
	if len(mem.store) >= mem.Capacity {
		delCount := mem.expireClean()
		if delCount <= 0 {
			return ErrKeysOverCapacity
		}
	}
	return nil
}

// 在超量后，才执行此函数
func (mem *MemCache) expireClean() int {
	delCount := 0
	for k, v := range mem.store {
		val := v.(ttlValue)
		if val.Expired() {
			delete(mem.store, k)
			delCount++
		}
	}
	return delCount
}

// autoExpireClean 自动在一定时间内清理过期key
// 当设置了大量的 expire key 且通常只读取一次的情况下再建议使用。
// interval 建议设置大一点，否则可能影响读性能，建议设置 5 minute
// go AutoCleanExpireKey(5 * time.Minute)
func (mem *MemCache) autoExpireClean(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			mem.rwMutex.Lock()
			mem.expireClean()
			mem.rwMutex.Unlock()
		}
	}
}

func (mem *MemCache) Save() *StatusCmd {
	if mem.Filename == "" {
		return &StatusCmd{baseCmd: baseCmd{err: ErrFilenameEmpty}}
	}

	disk, err := NewDisk(mem.Filename)
	if err != nil {
		return &StatusCmd{baseCmd: baseCmd{err: err}}
	}
	mem.rwMutex.RLock()
	defer mem.rwMutex.RUnlock()
	byt, err := json.Marshal(mem.store)
	if err != nil {
		return &StatusCmd{baseCmd: baseCmd{err: err}}
	}
	err = disk.WriteToFile(byt)
	return &StatusCmd{baseCmd: baseCmd{err: err}}
}

func (mem *MemCache) load() error {
	if mem.Filename == "" {
		return nil
	}
	values := make(map[string]ttlValue, 0)
	disk, err := NewDisk(mem.Filename)
	if err != nil {
		return err
	}
	byt, err := disk.ReadFromFile()
	if err != nil {
		return err
	}

	if len(byt) == 0 {
		return nil
	}
	if err := json.Unmarshal(byt, &values); err != nil {
		return err
	}
	mem.rwMutex.Lock()
	defer mem.rwMutex.Unlock()

	for k, v := range values {
		if !v.Expired() {
			mem.store[k] = v
		}
	}
	return nil
}

type Disk struct {
	filename string
}

//  NewDisk  filename
func NewDisk(filename string) (*Disk, error) {

	f, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	filename = f

	d := Disk{
		filename: filename,
	}
	return &d, nil
}

// WriteToFile 如果之前文件存在，则删除
func (d *Disk) WriteToFile(data []byte) error {
	if FilenameExists(d.filename) {
		if err := os.Remove(d.filename); err != nil {
			return err
		}
	} else {
		dir := filepath.Dir(d.filename)
		if err := os.MkdirAll(dir, os.FileMode(0666)); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(d.filename, os.O_RDWR|os.O_CREATE, os.FileMode(0666))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// ReadFromFile
func (d *Disk) ReadFromFile() ([]byte, error) {
	data := make([]byte, 0)
	if !FilenameExists(d.filename) {
		return data, nil
	}

	file, err := os.Open(d.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err = ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetCurrentDir() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		return "", err
	}
	return strings.Replace(dir, "\\", "/", -1), nil //将\替换成/
}

func FilenameExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}
