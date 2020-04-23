package local_cache

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

type sourceFunc = func(ctx context.Context) error

// author by suyue
// 本地缓存简单的实现
// 并不会删除已经过期 key 的值，因此只适用于 key 固定的情况，否则一定会出现内存泄漏
type store struct {
	ret         interface{}
	lastCalled  time.Time
	calledTimes int64
}

type localCache struct {
	m sync.Map
}

func (l *localCache) isExpire(s *store, e time.Duration, m int64) bool {
	return s.lastCalled.IsZero() || time.Now().After(s.lastCalled.Add(e)) || m <= atomic.LoadInt64(&s.calledTimes)
}

// key： 唯一的 namespace，同一个 key 会返回相同的值
// retAddr：保存结果的变量，必须是指针
// source Func：当被击穿后回源的函数。这个函数内部需要向 retAddr 赋值
// expire： 被击穿前最大的时间
// maxTimes： 被击穿前最大的次数
func (l *localCache) Load(
	ctx context.Context,
	key string,
	retAddr interface{},
	sf sourceFunc,
	expire time.Duration,
	maxTimes int64,
) (err error) {
	rv := reflect.ValueOf(retAddr)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("retAddr must be a pointer")
	}
	if value, ok := l.m.Load(key); ok {
		v, ok := value.(*store)
		if !ok {
			err = errors.New("wrong type of store")
			return
		}
		if !l.isExpire(v, expire, maxTimes) {
			atomic.AddInt64(&v.calledTimes, 1)
			if l.setValue(retAddr, v.ret) {
				return
			}
		}
	}
	if err := sf(ctx); err == nil {
		sr := &store{
			ret:        retAddr,
			lastCalled: time.Now(),
		}
		l.m.Store(key, sr)
	}
	return
}

func (l *localCache) setValue(dst, src interface{}) (ok bool) {
	ok = true
	func() {
		if err := recover(); err != nil {
			ok = false
		}
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(src).Elem())
	}()
	return
}

var DefaultLocalCache *localCache

func NewLocalCache() *localCache {
	var m sync.Map
	return &localCache{
		m: m,
	}
}

func init() {
	DefaultLocalCache = NewLocalCache()
}
