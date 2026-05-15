package testutils

import (
	"time"

	"github.com/stretchr/testify/mock"
)

// MockRedisClient は RedisClient インターフェースのモックです。
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) SetURLRecord(id string, baseUrl string, isSandCushion bool, isPublicCtrl bool, expireDelta *time.Duration) error {
	args := m.Called(id, baseUrl, isSandCushion, isPublicCtrl, expireDelta)
	return args.Error(0)
}

func (m *MockRedisClient) GetBaseUrl(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisClient) GetIsNeedCusionPage(key string) (bool, error) {
	args := m.Called(key)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisClient) GetIsPublicCtrl(key string) (bool, error) {
	args := m.Called(key)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisClient) IsExists(key string) (bool, error) {
	args := m.Called(key)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisClient) Close() error {
	args := m.Called()
	return args.Error(0)
}
