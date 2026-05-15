package redisclient

import (
	"fmt"
	"strconv"
)

func (r *RedisAdapter) GetBaseUrl(key string) (string, error) {
	baseUrl, err := r.Client.HGet(key, "base_url").Result()
	if err != nil {
		return "", err
	}

	return baseUrl, nil
}

func (r *RedisAdapter) GetIsNeedCusionPage(key string) (bool, error) {
	redisVal, err := r.Client.HGet(key, "cushion").Result()
	if err != nil {
		return false, fmt.Errorf("get redis: %w", err)
	}

	isNeed, err := strconv.ParseBool(redisVal)
	if err != nil {
		return false, fmt.Errorf("parse got val: %w", err)
	}

	return isNeed, nil
}

func (r *RedisAdapter) GetIsPublicCtrl(key string) (bool, error) {
	redisVal, err := r.Client.HGet(key, "public_ctrl").Result()
	if err != nil {
		return false, fmt.Errorf("get redis: %w", err)
	}

	isPublic, err := strconv.ParseBool(redisVal)
	if err != nil {
		return false, fmt.Errorf("parse got val: %w", err)
	}

	return isPublic, nil
}
