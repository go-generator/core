package template

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrNotFoundInCache = errors.New("not found in cache")
)

type cachedTemplate struct {
	Template
	expireAtTimestamp int64
}

type ILocalCache interface {
	Save(id string, t Template, expireAtTimestamp int64)
	Read(id string) (Template, error)
	Delete(id string)
}

type localCache struct {
	stop chan struct{}

	wg        sync.WaitGroup
	mu        sync.RWMutex
	Templates map[string]cachedTemplate
}

func NewLocalCache(cleanupInterval time.Duration) *localCache {
	lc := &localCache{
		Templates: make(map[string]cachedTemplate),
		stop:      make(chan struct{}),
	}
	if cleanupInterval.Seconds() != 0 {
		lc.wg.Add(1)
		go func(cleanupInterval time.Duration) {
			defer lc.wg.Done()
			lc.CleanupLoop(cleanupInterval)
		}(cleanupInterval)

	}
	return lc
}

func (lc *localCache) CleanupLoop(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-lc.stop:
			return
		case <-t.C:
			lc.mu.Lock()
			for uid, cu := range lc.Templates {
				if cu.expireAtTimestamp <= time.Now().Unix() {
					delete(lc.Templates, uid)
				}
			}
			lc.mu.Unlock()
		}
	}
}

func (lc *localCache) StopCleanup() {
	close(lc.stop)
	lc.wg.Wait()
}

func (lc *localCache) Save(id string, u Template, expireAtTimestamp int64) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.Templates[id] = cachedTemplate{
		Template:          u,
		expireAtTimestamp: expireAtTimestamp,
	}
}

func (lc *localCache) Read(id string) (Template, error) {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	cu, ok := lc.Templates[id]
	if !ok {
		return Template{}, ErrNotFoundInCache
	}

	return cu.Template, nil
}

func (lc *localCache) Delete(id string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	delete(lc.Templates, id)
}
