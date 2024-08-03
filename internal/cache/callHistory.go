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
	callHistoryCachePrefixKey = "callHistory:"
	// CallHistoryExpireTime expire time
	CallHistoryExpireTime = 5 * time.Minute
)

var _ CallHistoryCache = (*callHistoryCache)(nil)

// CallHistoryCache cache interface
type CallHistoryCache interface {
	Set(ctx context.Context, id uint64, data *model.CallHistory, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.CallHistory, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.CallHistory, error)
	MultiSet(ctx context.Context, data []*model.CallHistory, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// callHistoryCache define a cache struct
type callHistoryCache struct {
	cache cache.Cache
}

// NewCallHistoryCache new a cache
func NewCallHistoryCache(cacheType *model.CacheType) CallHistoryCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.CallHistory{}
		})
		return &callHistoryCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.CallHistory{}
		})
		return &callHistoryCache{cache: c}
	}

	return nil // no cache
}

// GetCallHistoryCacheKey cache key
func (c *callHistoryCache) GetCallHistoryCacheKey(id uint64) string {
	return callHistoryCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *callHistoryCache) Set(ctx context.Context, id uint64, data *model.CallHistory, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetCallHistoryCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *callHistoryCache) Get(ctx context.Context, id uint64) (*model.CallHistory, error) {
	var data *model.CallHistory
	cacheKey := c.GetCallHistoryCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *callHistoryCache) MultiSet(ctx context.Context, data []*model.CallHistory, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetCallHistoryCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *callHistoryCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.CallHistory, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetCallHistoryCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.CallHistory)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.CallHistory)
	for _, id := range ids {
		val, ok := itemMap[c.GetCallHistoryCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *callHistoryCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetCallHistoryCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *callHistoryCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetCallHistoryCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
