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
	distributionCachePrefixKey = "distribution:"
	// DistributionExpireTime expire time
	DistributionExpireTime = 5 * time.Minute
)

var _ DistributionCache = (*distributionCache)(nil)

// DistributionCache cache interface
type DistributionCache interface {
	Set(ctx context.Context, id uint64, data *model.Distribution, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Distribution, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Distribution, error)
	MultiSet(ctx context.Context, data []*model.Distribution, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// distributionCache define a cache struct
type distributionCache struct {
	cache cache.Cache
}

// NewDistributionCache new a cache
func NewDistributionCache(cacheType *model.CacheType) DistributionCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Distribution{}
		})
		return &distributionCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Distribution{}
		})
		return &distributionCache{cache: c}
	}

	return nil // no cache
}

// GetDistributionCacheKey cache key
func (c *distributionCache) GetDistributionCacheKey(id uint64) string {
	return distributionCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *distributionCache) Set(ctx context.Context, id uint64, data *model.Distribution, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetDistributionCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *distributionCache) Get(ctx context.Context, id uint64) (*model.Distribution, error) {
	var data *model.Distribution
	cacheKey := c.GetDistributionCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *distributionCache) MultiSet(ctx context.Context, data []*model.Distribution, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetDistributionCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *distributionCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Distribution, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetDistributionCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Distribution)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Distribution)
	for _, id := range ids {
		val, ok := itemMap[c.GetDistributionCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *distributionCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetDistributionCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *distributionCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetDistributionCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
