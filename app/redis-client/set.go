package redisclient

import (
	"fmt"
	"time"
)

func (r *RedisAdapter) SetURLRecord(id string, baseUrl string, isSandCushion bool, isPublicCtrl bool, expireDelta *time.Duration) error {
	// RedisにURLレコードを保存
	if err := r.Client.HMSet(id, map[string]interface{}{
		"base_url":    baseUrl,
		"cushion":     isSandCushion,
		"public_ctrl": isPublicCtrl,
	}).Err(); err != nil {
		return fmt.Errorf("setRecord: %w", err)
	}

	// 有効期限が指定されている場合、設定
	if expireDelta != nil {
		if err := r.Client.Expire(id, *expireDelta).Err(); err != nil {
			return fmt.Errorf("setExpire: %w", err)
		}
	}

	return nil
}
