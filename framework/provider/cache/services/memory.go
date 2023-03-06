package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/awaketai/goweb/framework"
	"strings"
	"sync"
	"time"
)

type MemoryData struct {
	val        any
	createTime time.Time
	ttl        time.Duration
}

type MemoryCache struct {
	container framework.Container
	data      map[string]*MemoryData
	lock      sync.RWMutex
}

func NewMemoryCache(params ...any) (any, error) {
	container := params[0].(framework.Container)
	obj := &MemoryCache{
		container: container,
		data:      map[string]*MemoryData{},
		lock:      sync.RWMutex{},
	}
	return obj, nil
}

func (m *MemoryCache) Get(ctx context.Context, key string) (string, error) {
	var val string
	if err := m.GetObj(ctx, key, &val); err != nil {
		return "", nil
	}
	return val, nil
}

func (m *MemoryCache) GetObj(ctx context.Context, key string, obj any) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if md, ok := m.data[key]; ok {
		if md.ttl != NoneDuration && time.Now().Sub(md.createTime) > md.ttl {
			delete(m.data, key)
			return ErrKeyNotFound
		}
		bt, err := json.Marshal(md.val)
		if err != nil {
			return err
		}
		err = json.Unmarshal(bt, obj)
		if err != nil {
			return err
		}
		return nil
	}
	return ErrKeyNotFound
}

func (m *MemoryCache) GetMany(ctx context.Context, keys []string) (map[string]string, error) {
	errs := make([]string, 0, len(keys))
	rets := make(map[string]string)
	for _, key := range keys {
		val, err := m.Get(ctx, key)
		if err != nil {
			errs = append(errs, "key:"+key+" "+err.Error())
			continue
		}
		rets[key] = val
	}
	if len(errs) == 0 {
		return rets, nil
	}
	return rets, errors.New(strings.Join(errs, "||"))
}

func (m *MemoryCache) Set(ctx context.Context, key string, val any, timeout time.Duration) error {
	return m.SetObj(ctx, key, val, timeout)
}

func (m *MemoryCache) SetObj(ctx context.Context, key string, val any, timeout time.Duration) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	md := &MemoryData{
		val:        val,
		createTime: time.Now(),
		ttl:        timeout,
	}
	m.data[key] = md
	return nil
}

func (m *MemoryCache) SetMany(ctx context.Context, data map[string]string, timeout time.Duration) error {
	var errs []string
	for k, v := range data {
		err := m.Set(ctx, k, v, timeout)
		if err != nil {
			errs = append(errs, "key:"+k+" "+err.Error())
			continue
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "||"))
	}
	return nil
}

func (m *MemoryCache) SetTTL(ctx context.Context, key string, timeout time.Duration) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if md, ok := m.data[key]; ok {
		md.ttl = timeout
		return nil
	}
	return ErrKeyNotFound
}

func (m *MemoryCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if md, ok := m.data[key]; ok {
		return md.ttl, nil
	}
	return NoneDuration, ErrKeyNotFound

}

func (m *MemoryCache) Del(ctx context.Context, key string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.data, key)
	return nil
}

func (m *MemoryCache) DelMany(ctx context.Context, keys []string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, key := range keys {
		delete(m.data, key)
	}
	return nil
}
