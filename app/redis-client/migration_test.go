package redisclient_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrateAddPublicCtrl(t *testing.T) {
	mr, adapter, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// 1. 既存データ（public_ctrl なし）
	id1 := "mig_id1"
	mr.HSet(id1, "base_url", "https://example.com/1")
	mr.HSet(id1, "cushion", "true")

	// 2. 既存データ（既に public_ctrl あり）
	id2 := "mig_id2"
	mr.HSet(id2, "base_url", "https://example.com/2")
	mr.HSet(id2, "public_ctrl", "false")

	// 3. ハッシュ型ではないキー（無視されるべき）
	id3 := "mig_id3"
	mr.Set(id3, "just-a-string")

	// マイグレーション実行
	err := adapter.MigrateAddPublicCtrl(true)
	assert.NoError(t, err)

	// 検証: id1 にはデフォルト値 true がセットされていること
	val1, _ := strconv.ParseBool(mr.HGet(id1, "public_ctrl"))
	assert.Equal(t, true, val1)

	// 検証: id2 の既存値 false は上書きされていないこと
	val2, _ := strconv.ParseBool(mr.HGet(id2, "public_ctrl"))
	assert.Equal(t, false, val2)

	// 検証: id3 はハッシュ型ではないので public_ctrl は存在しないはず
	assert.Empty(t, mr.HGet(id3, "public_ctrl"))
}
