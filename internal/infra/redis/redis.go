package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"project-api/internal/infra/config"
	"project-api/internal/infra/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisClient struct {
	*redis.Client
	TTL time.Duration // Default TTL
}

func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:         config.Config.Redis.Endpoint,
		Password:     config.Config.Redis.Password,
		DB:           0,
		MaxRetries:   3,
		PoolSize:     10,
		MinIdleConns: 2,
	})

	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}
	return &RedisClient{
		Client: client,
		TTL:    24 * time.Hour, // Default TTL
	}
}

// SetToCache เก็บข้อมูลใดๆ ลง Redis ด้วย key ที่กำหนด
func (r *RedisClient) SetToCache(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		logger.Error("Failed to marshal value for cache", zap.String("key", key), zap.Error(err))
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	if err := r.Client.SetEx(ctx, key, data, r.TTL).Err(); err != nil { // เปลี่ยน SetEX เป็น SetEx
		logger.Error("Failed to set value to Redis", zap.String("key", key), zap.Error(err))
		return fmt.Errorf("failed to set value to cache: %w", err)
	}
	return nil
}

// GetFromCache ดึงข้อมูลจาก Redis ด้วย key ที่กำหนด และคืนเป็น []byte เพื่อให้ caller unmarshal เอง
func (r *RedisClient) GetFromCache(ctx context.Context, key string) ([]byte, error) {
	data, err := r.Client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		logger.Error("Failed to get data from Redis", zap.String("key", key), zap.Error(err))
		return nil, fmt.Errorf("failed to get data from cache: %w", err)
	}
	return data, nil
}

// DeleteFromCache ลบ cache ด้วย key ที่กำหนด
func (r *RedisClient) DeleteFromCache(ctx context.Context, key string) error {
	if err := r.Client.Del(ctx, key).Err(); err != nil {
		logger.Warn("Failed to delete cache", zap.String("key", key), zap.Error(err))
		return fmt.Errorf("failed to delete cache: %w", err)
	}
	return nil
}

// BatchSetToCache เก็บข้อมูลหลายรายการลง Redis พร้อมกัน
func (r *RedisClient) BatchSetToCache(ctx context.Context, values map[string]interface{}) error {
	pipe := r.Client.Pipeline()
	for key, value := range values {
		data, err := json.Marshal(value)
		if err != nil {
			logger.Error("Failed to marshal value for batch set", zap.String("key", key), zap.Error(err))
			continue
		}
		pipe.SetEx(ctx, key, data, r.TTL) // เปลี่ยน SetEX เป็น SetEx
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		logger.Error("Failed to execute batch set in Redis", zap.Error(err))
		return fmt.Errorf("failed to batch set values: %w", err)
	}
	return nil
}

// BatchGetFromCache ดึงข้อมูลหลายรายการจาก Redis พร้อมกัน และคืนเป็น map[string][]byte
func (r *RedisClient) BatchGetFromCache(ctx context.Context, keys []string) (map[string][]byte, error) {
	values, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		logger.Error("Failed to batch get data from Redis", zap.Error(err))
		return nil, fmt.Errorf("failed to batch get data: %w", err)
	}

	result := make(map[string][]byte)
	for i, val := range values {
		if val == nil {
			continue // Cache miss สำหรับ key นี้
		}
		result[keys[i]] = []byte(val.(string))
	}
	return result, nil
}

// InvalidateCache ลบ cache หลาย key พร้อมกัน
func (r *RedisClient) InvalidateCache(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil // ไม่มี key ให้ลบ
	}
	if err := r.Client.Del(ctx, keys...).Err(); err != nil {
		logger.Warn("Failed to invalidate cache", zap.Strings("keys", keys), zap.Error(err))
		return fmt.Errorf("failed to invalidate cache: %w", err)
	}
	return nil
}
