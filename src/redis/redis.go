package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	conn *redis.Client
}

func (r *Redis) Connect(uri string, database int) error {
	c := redis.NewClient(&redis.Options{
		Addr: uri,
		DB:   database,
	})

	r.conn = c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	return c.Ping(ctx).Err()
}

func (r *Redis) GetString(key string) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	existsResult := r.conn.Exists(ctx, key)

	if err := existsResult.Err(); err != nil {
		return "", false, err
	}

	if existsResult.Val() == 0 {
		return "", false, nil
	}

	result := r.conn.Get(ctx, key)

	return result.Val(), true, result.Err()
}

func (r *Redis) GetBytes(key string) ([]byte, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	existsResult := r.conn.Exists(ctx, key)

	if err := existsResult.Err(); err != nil {
		return nil, false, err
	}

	if existsResult.Val() == 0 {
		return nil, false, nil
	}

	result := r.conn.Get(ctx, key)

	if err := result.Err(); err != nil {
		return nil, true, err
	}

	data, err := result.Bytes()

	return data, true, err
}

func (r *Redis) Exists(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	result := r.conn.Exists(ctx, key)

	return result.Val() == 1, result.Err()
}

func (r *Redis) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	return r.conn.Del(ctx, key).Err()
}

func (r *Redis) Set(key string, value interface{}, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	return r.conn.Set(ctx, key, value, ttl).Err()
}

func (r *Redis) Close() error {
	return r.conn.Close()
}
