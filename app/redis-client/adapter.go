package redisclient

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type RedisClient interface {
	SetURLRecord(id string, baseUrl string, isSandCushion bool, isPublicCtrl bool, expireDelta *time.Duration) error
	GetBaseUrl(key string) (string, error)
	GetIsNeedCusionPage(key string) (bool, error)
	GetIsPublicCtrl(key string) (bool, error)
	IsExists(key string) (bool, error)
	Close() error
}

type RedisAdapter struct {
	Client *redis.Client
}

func NewRedisAdapter(addr string, password string, db int) (*RedisAdapter, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, fmt.Errorf("connectRedis: %w", err)
	}

	return &RedisAdapter{Client: client}, nil
}

func (r *RedisAdapter) Set(key string, value string) error {
	return r.Client.Set(key, value, 0).Err()
}

func (r *RedisAdapter) Get(key string) (string, error) {
	return r.Client.Get(key).Result()
}

func (r *RedisAdapter) Close() error {
	return r.Client.Close()
}

func (r *RedisAdapter) IsExists(key string) (bool, error) {
	exists, err := r.Client.Exists(key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}
