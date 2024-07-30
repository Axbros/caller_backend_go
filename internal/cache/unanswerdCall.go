package cache

import (
	"context"
	"strings"
	"time"

	"github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/encoding"
	"github.com/zhufuyi/sponge/pkg/utils"

	"caller/internal/model"
)

const (
	// cache prefix key, must end with a colon
	unanswerdCallCachePrefixKey = "unanswerdCall:"
	// UnanswerdCallExpireTime expire time
	UnanswerdCallExpireTime = 5 * time.Minute
)

var _ UnanswerdCallCache = (*unanswerdCallCache)(nil)

// UnanswerdCallCache cache interface
type UnanswerdCallCache interface {
	Set(ctx context.Context, id uint64, data *model.UnanswerdCall, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.UnanswerdCall, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.UnanswerdCall, error)
	MultiSet(ctx context.Context, data []*model.UnanswerdCall, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// unanswerdCallCache define a cache struct
type unanswerdCallCache struct {
	cache cache.Cache
}

// NewUnanswerdCallCache new a cache
func NewUnanswerdCallCache(cacheType *model.CacheType) UnanswerdCallCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.UnanswerdCall{}
		})
		return &unanswerdCallCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.UnanswerdCall{}
		})
		return &unanswerdCallCache{cache: c}
	}

	return nil // no cache
}

// GetUnanswerdCallCacheKey cache key
func (c *unanswerdCallCache) GetUnanswerdCallCacheKey(id uint64) string {
	return unanswerdCallCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *unanswerdCallCache) Set(ctx context.Context, id uint64, data *model.UnanswerdCall, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetUnanswerdCallCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *unanswerdCallCache) Get(ctx context.Context, id uint64) (*model.UnanswerdCall, error) {
	var data *model.UnanswerdCall
	cacheKey := c.GetUnanswerdCallCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *unanswerdCallCache) MultiSet(ctx context.Context, data []*model.UnanswerdCall, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetUnanswerdCallCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *unanswerdCallCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.UnanswerdCall, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetUnanswerdCallCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.UnanswerdCall)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.UnanswerdCall)
	for _, id := range ids {
		val, ok := itemMap[c.GetUnanswerdCallCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *unanswerdCallCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetUnanswerdCallCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *unanswerdCallCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetUnanswerdCallCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
