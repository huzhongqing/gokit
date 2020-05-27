package redislock

import (
	"fmt"
	"testing"
	"time"
)

func TestClientNew(t *testing.T) {
	lockKey := "Locker_Test"
	lockCli := New("192.168.8.145:6379", "12345678")
	locker, err := lockCli.Obtain(lockKey, 10*time.Second)
	if err == ErrNotObtained {
		t.Fatal("锁必需获取的到")
	} else if err != nil {
		t.Fatal(err)
	}
	defer locker.Release()

	go func() {
		time.AfterFunc(time.Second, func() {
			fmt.Println("goroutine 去竞争锁")
			_, err := lockCli.Obtain(lockKey, 10*time.Second)
			if err != ErrNotObtained {
				t.Fatal("锁应该是获取不到的")
			}
		})
	}()

	ttl, err := locker.TTL()
	if err != nil {
		t.Fatal(err)
	}
	if ttl > time.Second {
		fmt.Println("锁有效")
	}
	time.Sleep(3 * time.Second)

	if err := locker.Refresh(10 * time.Second); err != nil {
		t.Fatal(err)
	}

	ttl, err = locker.TTL()
	if err != nil {
		t.Fatal(err)
	}
	if ttl < 8 {
		t.Fatal("Refresh 失效")
	}
}
