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
	clientsCachePrefixKey = "clients:"
	// ClientsExpireTime expire time
	ClientsExpireTime = 5 * time.Minute
)

var _ ClientsCache = (*clientsCache)(nil)

// ClientsCache cache interface
type ClientsCache interface {
	Set(ctx context.Context, id uint64, data *model.Clients, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Clients, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Clients, error)
	MultiSet(ctx context.Context, data []*model.Clients, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// clientsCache define a cache struct
type clientsCache struct {
	cache cache.Cache
}

// NewClientsCache new a cache
func NewClientsCache(cacheType *model.CacheType) ClientsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Clients{}
		})
		return &clientsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Clients{}
		})
		return &clientsCache{cache: c}
	}

	return nil // no cache
}

// GetClientsCacheKey cache key
func (c *clientsCache) GetClientsCacheKey(id uint64) string {
	return clientsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *clientsCache) Set(ctx context.Context, id uint64, data *model.Clients, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetClientsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *clientsCache) Get(ctx context.Context, id uint64) (*model.Clients, error) {
	var data *model.Clients
	cacheKey := c.GetClientsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *clientsCache) MultiSet(ctx context.Context, data []*model.Clients, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetClientsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *clientsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Clients, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetClientsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Clients)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Clients)
	for _, id := range ids {
		val, ok := itemMap[c.GetClientsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *clientsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetClientsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *clientsCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetClientsCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
