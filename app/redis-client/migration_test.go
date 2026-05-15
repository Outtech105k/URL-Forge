package redisclient_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrateAddPublicCtrl(t *testing.T) {
	tests := []struct {
		name         string
		defaultValue bool
	}{
		{"DefaultTrue", true},
		{"DefaultFalse", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr, adapter, cleanup := setupTestEnvironment(t)
			defer cleanup()

			// 1. 既存データ（public_ctrl なし）
			id1 := "mig_id1"
			mr.HSet(id1, "base_url", "https://example.com/1")
			mr.HSet(id1, "cushion", "true")

			// 2. 既存データ（既に public_ctrl あり）
			// デフォルト値と逆の値をセットして、上書きされないことを確認する
			id2 := "mig_id2"
			existingVal := !tt.defaultValue
			mr.HSet(id2, "base_url", "https://example.com/2")
			mr.HSet(id2, "public_ctrl", strconv.FormatBool(existingVal))

			// 3. ハッシュ型ではないキー（無視されるべき）
			id3 := "mig_id3"
			mr.Set(id3, "just-a-string")

			// マイグレーション実行
			err := adapter.MigrateAddPublicCtrl(tt.defaultValue)
			assert.NoError(t, err)

			// 検証: id1 には指定したデフォルト値がセットされていること
			val1, _ := strconv.ParseBool(mr.HGet(id1, "public_ctrl"))
			assert.Equal(t, tt.defaultValue, val1)

			// 検証: id2 の既存値は上書きされていないこと
			val2, _ := strconv.ParseBool(mr.HGet(id2, "public_ctrl"))
			assert.Equal(t, existingVal, val2)

			// 検証: id3 はハッシュ型ではないので public_ctrl は存在しないはず
			assert.Empty(t, mr.HGet(id3, "public_ctrl"))
		})
	}
}
