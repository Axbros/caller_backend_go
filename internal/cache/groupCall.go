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
	groupCallCachePrefixKey = "groupCall:"
	// GroupCallExpireTime expire time
	GroupCallExpireTime = 5 * time.Minute
)

var _ GroupCallCache = (*groupCallCache)(nil)

// GroupCallCache cache interface
type GroupCallCache interface {
	Set(ctx context.Context, id uint64, data *model.GroupCall, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.GroupCall, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.GroupCall, error)
	MultiSet(ctx context.Context, data []*model.GroupCall, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// groupCallCache define a cache struct
type groupCallCache struct {
	cache cache.Cache
}

// NewGroupCallCache new a cache
func NewGroupCallCache(cacheType *model.CacheType) GroupCallCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.GroupCall{}
		})
		return &groupCallCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.GroupCall{}
		})
		return &groupCallCache{cache: c}
	}

	return nil // no cache
}

// GetGroupCallCacheKey cache key
func (c *groupCallCache) GetGroupCallCacheKey(id uint64) string {
	return groupCallCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *groupCallCache) Set(ctx context.Context, id uint64, data *model.GroupCall, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetGroupCallCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *groupCallCache) Get(ctx context.Context, id uint64) (*model.GroupCall, error) {
	var data *model.GroupCall
	cacheKey := c.GetGroupCallCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *groupCallCache) MultiSet(ctx context.Context, data []*model.GroupCall, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetGroupCallCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *groupCallCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.GroupCall, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetGroupCallCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.GroupCall)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.GroupCall)
	for _, id := range ids {
		val, ok := itemMap[c.GetGroupCallCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *groupCallCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetGroupCallCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *groupCallCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetGroupCallCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
