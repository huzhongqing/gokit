package redislock

import (
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v7"
)

var (
	// ErrNotObtained is returned when a lock cannot be obtained.
	ErrNotObtained = redislock.ErrNotObtained

	// ErrLockNotHeld is returned when trying to release an inactive lock.
	ErrLockNotHeld = redislock.ErrLockNotHeld
)

type Client interface {
	Obtain(key string, ttl time.Duration) (Locker, error)
}

type Locker interface {
	Key() string
	TTL() (time.Duration, error)
	Refresh(ttl time.Duration) error
	Release() error
}

func New(addr, password string) Client {
	return newIClient(addr, password)
}

func newIClient(addr, password string) *iClient {
	rCli := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     addr,
		Password: password,
	})
	rLock := redislock.New(rCli)
	return &iClient{
		lockClient: rLock,
	}
}

type iLock struct {
	*redislock.Lock
}

func (l *iLock) Refresh(ttl time.Duration) error {
	return l.Lock.Refresh(ttl, nil)
}

type iClient struct {
	lockClient *redislock.Client
	lock       *iLock
}

func (c *iClient) Obtain(key string, ttl time.Duration) (Locker, error) {
	l, err := c.lockClient.Obtain(key, ttl, nil)
	if err != nil {
		return nil, err
	}
	c.lock = &iLock{l}
	return c.lock, nil
}
