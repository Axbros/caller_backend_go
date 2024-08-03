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
	groupClientCachePrefixKey = "groupClient:"
	// GroupClientExpireTime expire time
	GroupClientExpireTime = 5 * time.Minute
)

var _ GroupClientCache = (*groupClientCache)(nil)

// GroupClientCache cache interface
type GroupClientCache interface {
	Set(ctx context.Context, id uint64, data *model.GroupClient, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.GroupClient, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.GroupClient, error)
	MultiSet(ctx context.Context, data []*model.GroupClient, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// groupClientCache define a cache struct
type groupClientCache struct {
	cache cache.Cache
}

// NewGroupClientCache new a cache
func NewGroupClientCache(cacheType *model.CacheType) GroupClientCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.GroupClient{}
		})
		return &groupClientCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.GroupClient{}
		})
		return &groupClientCache{cache: c}
	}

	return nil // no cache
}

// GetGroupClientCacheKey cache key
func (c *groupClientCache) GetGroupClientCacheKey(id uint64) string {
	return groupClientCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *groupClientCache) Set(ctx context.Context, id uint64, data *model.GroupClient, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetGroupClientCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *groupClientCache) Get(ctx context.Context, id uint64) (*model.GroupClient, error) {
	var data *model.GroupClient
	cacheKey := c.GetGroupClientCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *groupClientCache) MultiSet(ctx context.Context, data []*model.GroupClient, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetGroupClientCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *groupClientCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.GroupClient, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetGroupClientCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.GroupClient)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.GroupClient)
	for _, id := range ids {
		val, ok := itemMap[c.GetGroupClientCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *groupClientCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetGroupClientCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *groupClientCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetGroupClientCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
