package cache

import (
	"context"
	"errors"
	"sync"
	"time"
)

type LocalCache struct {
	data map[string]*cacheItem
	mu   sync.RWMutex
	max  int
	ttl  time.Duration
}

type cacheItem struct {
	value    interface{}
	expireAt time.Time
}

func NewLocalCache(max int, ttl time.Duration) *LocalCache {
	lc := &LocalCache{
		data: make(map[string]*cacheItem),
		max:  max,
		ttl:  ttl,
	}
	go lc.gc()
	return lc
}

func (l *LocalCache) Get(ctx context.Context, key string) (interface{}, error) {
	l.mu.RLock()
	item, ok := l.data[key]
	l.mu.RUnlock()

	if !ok || time.Now().After(item.expireAt) {
		return nil, errors.New("key not found")
	}
	return item.value, nil
}

func (l *LocalCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.max > 0 && len(l.data) >= l.max {
		l.evict()
	}

	l.data[key] = &cacheItem{
		value:    value,
		expireAt: time.Now().Add(ttl),
	}
	return nil
}

func (l *LocalCache) Delete(ctx context.Context, key string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.data, key)
	return nil
}

func (l *LocalCache) evict() {
	for k := range l.data {
		delete(l.data, k)
		break
	}
}

func (l *LocalCache) gc() {
	for {
		time.Sleep(time.Minute)
		l.mu.Lock()
		for k, v := range l.data {
			if time.Now().After(v.expireAt) {
				delete(l.data, k)
			}
		}
		l.mu.Unlock()
	}
}
