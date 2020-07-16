package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

// CacheOptions is config redis
type CacheOptions struct {
	Host          string
	Port          int
	Password      string
	MaxIdle       int
	MaxActive     int
	IdleTimeout   int
	Enabled       bool
	DatabaseIndex int
}

var pool *redis.Pool

// Connect is connection to redis
func Connect(cacheOptions CacheOptions) *redis.Pool {
	if pool == nil {
		pool = &redis.Pool{
			MaxIdle:     cacheOptions.MaxIdle,
			MaxActive:   cacheOptions.MaxActive,
			IdleTimeout: time.Duration(cacheOptions.IdleTimeout) * time.Second,
			Dial: func() (redis.Conn, error) {
				address := fmt.Sprintf("%s:%d", cacheOptions.Host, cacheOptions.Port)
				c, err := redis.Dial(
					"tcp",
					address,
					redis.DialDatabase(cacheOptions.DatabaseIndex),
					redis.DialPassword(cacheOptions.Password),
				)
				if err != nil {
					return nil, err
				}

				// Do authentication process if password not empty.
				if cacheOptions.Password != "" {
					if _, err := c.Do("AUTH", cacheOptions.Password); err != nil {
						c.Close()
						return nil, err
					}
				}

				return c, nil
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}

				_, err := c.Do("PING")
				return err
			},
			Wait:            true,
			MaxConnLifetime: 15 * time.Minute,
		}

		return pool
	}

	return pool
}

// Command General function
func Command(ctx context.Context, cachePool *redis.Pool, command, keyCache string, data interface{}) (interface{}, error) {
	// store rule to cache
	conn, err := cachePool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ruleByte, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	data, err = conn.Do(command, keyCache, ruleByte)
	if err != nil {
		return data, err
	}

	return data, nil
}

// Set save data to redis
func Set(ctx context.Context, cachePool *redis.Pool, keyCache string, data interface{}) error {
	// store rule to cache
	conn, err := cachePool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	ruleByte, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", keyCache, ruleByte)
	if err != nil {
		return err
	}

	return nil
}

// SetEx is save data to redis with time
func SetEx(ctx context.Context, cachePool *redis.Pool, keyCache string, expire int, data interface{}) error {
	// store rule to cache
	conn, err := cachePool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	ruleByte, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SETEX", keyCache, expire, ruleByte)
	if err != nil {
		return err
	}

	return nil
}

// Get is  function to get data from redis with key
func Get(ctx context.Context, cachePool *redis.Pool, keyCache string) ([]byte, error) {
	conn, err := cachePool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ruleByte, err := redis.Bytes(conn.Do("GET", keyCache))
	if err != nil {
		return nil, err
	}
	return ruleByte, nil
}

// Delete is function to remove data from redis by key
func Delete(ctx context.Context, cachePool *redis.Pool, keyCache string) error {
	conn, err := cachePool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("DEL", keyCache)
	if err != nil {
		return err
	}
	return nil
}

// HSet set item
func HSet(ctx context.Context, cachePool *redis.Pool, keyCache, field string, data interface{}) error {
	// store rule to cache
	conn, err := cachePool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	message, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("HSET", keyCache, field, message)
	if err != nil {
		return err
	}

	return nil
}

// HGetAll set item
func HGetAll(ctx context.Context, cachePool *redis.Pool, keyCache string) ([][]byte, error) {
	// store rule to cache
	conn, err := cachePool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	data, err := conn.Do("HGETALL", keyCache)
	if err != nil {
		return nil, err
	}

	arrayData, err := redis.ByteSlices(data, err)
	if err != nil {
		return nil, err
	}

	return arrayData, nil
}

// TTL is function to get time life from data
func TTL(ctx context.Context, cachePool *redis.Pool, keyCache string) (int, error) {
	conn, err := cachePool.GetContext(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	ttl, err := redis.Int(conn.Do("TTL", keyCache))
	if err != nil {
		return 0, err
	}
	return ttl, nil
}

// Expire is set data with time
func Expire(ctx context.Context, cachePool *redis.Pool, keyCache string, expire int) error {
	conn, err := cachePool.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("EXPIRE", keyCache, expire)
	if err != nil {
		return err
	}
	return nil
}
