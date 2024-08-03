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
	smsCachePrefixKey = "sms:"
	// SmsExpireTime expire time
	SmsExpireTime = 5 * time.Minute
)

var _ SmsCache = (*smsCache)(nil)

// SmsCache cache interface
type SmsCache interface {
	Set(ctx context.Context, id uint64, data *model.Sms, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Sms, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Sms, error)
	MultiSet(ctx context.Context, data []*model.Sms, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// smsCache define a cache struct
type smsCache struct {
	cache cache.Cache
}

// NewSmsCache new a cache
func NewSmsCache(cacheType *model.CacheType) SmsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Sms{}
		})
		return &smsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Sms{}
		})
		return &smsCache{cache: c}
	}

	return nil // no cache
}

// GetSmsCacheKey cache key
func (c *smsCache) GetSmsCacheKey(id uint64) string {
	return smsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *smsCache) Set(ctx context.Context, id uint64, data *model.Sms, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetSmsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *smsCache) Get(ctx context.Context, id uint64) (*model.Sms, error) {
	var data *model.Sms
	cacheKey := c.GetSmsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *smsCache) MultiSet(ctx context.Context, data []*model.Sms, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetSmsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *smsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Sms, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetSmsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Sms)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Sms)
	for _, id := range ids {
		val, ok := itemMap[c.GetSmsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *smsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetSmsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *smsCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetSmsCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
