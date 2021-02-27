package redis

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
