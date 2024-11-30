package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/zhufuyi/sponge/pkg/logger"
	"golang.org/x/sync/singleflight"
)

var _ RedisDao = (*redisDao)(nil)

// RedisDao defining the dao interface
type RedisDao interface {
	// SetIPAddrByMachineCode2WebsocketConnections(ctx context.Context, key string, value string) error
	GetIPAddrByMachineCodeFromWebsocketConnections(ctx context.Context, key string) (string, error)
	SetMessageStore(ctx context.Context, key string, value interface{}) error
	DeleteMessageStore(ctx context.Context, key string) error
	SetConn(ctx context.Context, key string, value interface{}) error
	GetConn(ctx context.Context, key string) (interface{}, error)
	SetKey(ctx context.Context, key string, value string) error
	GetKey(ctx context.Context, key string) (string, error)
	DeleteKey(ctx context.Context, key string) error
	GetQueenValue(ctx context.Context, key string) (string, error)

	PushClient2GroupName(ctx context.Context, prefix string, groupName string, clientID int) error
	DeleteRedisKeysWithPrefix(ctx context.Context, prefix string) error
	PushOfflinePhoneNumber(ctx context.Context, key string, value interface{}) error
	GetAllOfflinePhoneNumber(ctx context.Context, key string) ([]string, error)
}

type redisDao struct {
	client *redis.Client
	sfg    *singleflight.Group // if cache is nil, the sfg is not used.
}

// NewRedisDao creating the dao interface
func NewRedisDao(client *redis.Client) RedisDao {
	return &redisDao{
		client: client,
		sfg:    new(singleflight.Group),
	}
}

//	func (r *redisDao) SetIPAddrByMachineCode2WebsocketConnections(ctx context.Context, key string, value string) error {
//		// 设置 Hash 字段值
//		err := r.client.HSet(ctx, "websocket_connections", key, value).Err()
//		if err != nil {
//			logger.Errorf("设置 Hash 失败:", err)
//			return err
//		}
//		return nil
//	}
func (r *redisDao) PushOfflinePhoneNumber(ctx context.Context, key string, value interface{}) error {
	return r.client.LPush(ctx, key, value).Err()
}
func (r *redisDao) GetAllOfflinePhoneNumber(ctx context.Context, key string) ([]string, error) {
	return r.client.LRange(ctx, key, 0, -1).Result()
}

func (r *redisDao) PushClient2GroupName(ctx context.Context, prefix string, groupName string, clientID int) error {
	key := prefix + groupName
	return r.client.RPush(ctx, key, clientID).Err()
}
func (r *redisDao) DeleteRedisKeysWithPrefix(ctx context.Context, prefix string) error {
	// 获取所有以指定前缀开头的键
	keys, err := r.client.Keys(ctx, prefix+"*").Result()
	if err != nil {
		return errors.New("获取键列表失败:" + err.Error())
	}

	// 遍历键并判断类型，若为列表则删除
	for _, key := range keys {
		keyType, err := r.client.Type(ctx, key).Result()
		if err != nil {
			fmt.Println("获取键类型失败：", err)

			continue
		}
		if keyType == "list" {
			err = r.client.Del(ctx, key).Err()
			if err != nil {
				return errors.New("删除键失败" + err.Error())
			}
		}
	}
	return nil
}
func (r *redisDao) GetIPAddrByMachineCodeFromWebsocketConnections(ctx context.Context, key string) (string, error) {
	value, err := r.client.HGet(ctx, "websocket_connections", key).Result()
	if err == redis.Nil {
		logger.Error("键不存在")
		return "", nil
	} else if err != nil {
		logger.Errorf("获取值失败:", err)
		return "", err
	}
	return value, nil
}
func (r *redisDao) SetMessageStore(ctx context.Context, key string, value interface{}) error {
	// 设置 Hash 字段值
	err := r.client.HSet(ctx, "store", key, value).Err()
	if err != nil {
		logger.Errorf("设置 Store Message Hash 失败:", err)
		return err
	}
	return nil
}
func (r *redisDao) SetConn(ctx context.Context, key string, value interface{}) error {
	// 设置 Hash 字段值
	err := r.client.HSet(ctx, "connections", key, value).Err()
	if err != nil {
		logger.Errorf("Reids中设置 connections 失败:", err)
		return err
	}
	return nil
}
func (r *redisDao) GetConn(ctx context.Context, key string) (interface{}, error) {
	value, err := r.client.HGet(ctx, "store", key).Result()
	if err == redis.Nil {
		logger.Error("键不存在")
		return "", nil
	} else if err != nil {
		logger.Errorf("获取值失败:", err)
		return "", err
	}
	return value, nil
}
func (r *redisDao) DeleteMessageStore(ctx context.Context, key string) error {
	// 删除 Hash 字段
	err := r.client.HDel(ctx, "store", key).Err()
	if err != nil {
		logger.Errorf("删除 Store Message Hash 失败: %v", err)
		return err
	}
	logger.Info("delete store", logger.String("key", key))
	return nil
}
func (r *redisDao) SetKey(ctx context.Context, key string, value string) error {
	err := r.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set key in Redis: %v", err)
	}
	return nil
}

func (r *redisDao) GetKey(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", errors.New("key does not exist")
	} else if err != nil {
		return "", fmt.Errorf("failed to get key from Redis: %v", err)
	}
	return value, nil
}

func (r *redisDao) DeleteKey(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key from Redis: %v", err)
	}
	return nil
}
func (r *redisDao) GetQueenValue(ctx context.Context, key string) (string, error) {
	// 读取列表的首个元素
	value, err := r.client.LPop(ctx, key).Result()
	if err != nil {
		fmt.Println("读取列表元素时出错:", err)
		return "", err
	}
	// 将读取的元素添加到列表尾部
	err = r.client.RPush(ctx, key, value).Err()
	if err != nil {
		fmt.Println("将元素添加到列表尾部时出错:", err)
		return "", err
	}
	return value, nil
}
