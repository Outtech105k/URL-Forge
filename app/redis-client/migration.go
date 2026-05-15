package redisclient

import (
	"fmt"
)

func (r *RedisAdapter) MigrateAddPublicCtrl(defaultValue bool) error {
	iter := r.Client.Scan(0, "*", 0).Iterator()
	for iter.Next() {
		key := iter.Val()
		
		// キーがハッシュ型であることを確認
		t, err := r.Client.Type(key).Result()
		if err != nil {
			return fmt.Errorf("getType for key %s: %w", key, err)
		}
		
		if t == "hash" {
			// public_ctrl が存在しない場合のみデフォルト値をセット
			if err := r.Client.HSetNX(key, "public_ctrl", defaultValue).Err(); err != nil {
				return fmt.Errorf("hsetnx for key %s: %w", key, err)
			}
		}
	}
	
	if err := iter.Err(); err != nil {
		return fmt.Errorf("scan error: %w", err)
	}
	
	return nil
}
