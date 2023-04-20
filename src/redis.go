package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/draw"
	"time"

	"github.com/go-redis/redis/v8"
)

// Redis is a utility client for reading and writing values to the Redis server.
type Redis struct {
	conn *redis.Client
}

// Connect connects to the Redis server using the configuration values provided.
func (r *Redis) Connect(conf RedisConfig) error {
	c := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Username: conf.User,
		Password: conf.Password,
		DB:       conf.Database,
	})

	r.conn = c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	return c.Ping(ctx).Err()
}

// GetString gets the value from Redis by the key and returns the value as a string.
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

// GetBytes gets the value from Redis by the key and returns the value as a byte array.
func (r *Redis) GetBytes(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	existsResult := r.conn.Exists(ctx, key)

	if err := existsResult.Err(); err != nil {
		return nil, err
	}

	if existsResult.Val() == 0 {
		return nil, nil
	}

	result := r.conn.Get(ctx, key)

	if err := result.Err(); err != nil {
		return nil, err
	}

	data, err := result.Bytes()

	return data, err
}

// GetNRGBA gets the value from Redis by the key and returns the value as Go image.NRGBA image type.
func (r *Redis) GetNRGBA(key string) (*image.NRGBA, bool, error) {
	value, err := r.GetBytes(key)

	if err != nil {
		return nil, false, err
	}

	if value == nil {
		return nil, false, nil
	}

	img, format, err := image.Decode(bytes.NewReader(value))

	if err != nil {
		return nil, false, err
	}

	if format != "NRGBA" {
		outputImg := image.NewNRGBA(img.Bounds())

		draw.Draw(outputImg, img.Bounds(), img, image.Pt(0, 0), draw.Src)

		return outputImg, true, nil
	}

	return img.(*image.NRGBA), true, nil
}

// Exists returns whether or not the key exists in the database.
func (r *Redis) Exists(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	result := r.conn.Exists(ctx, key)

	return result.Val() == 1, result.Err()
}

// Delete deletes the key from the database if it exists.
func (r *Redis) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	return r.conn.Del(ctx, key).Err()
}

// Set puts the key-value into the database with an optional TTL value.
func (r *Redis) Set(key string, value interface{}, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	return r.conn.Set(ctx, key, value, ttl).Err()
}

// Close closes the connection to the database.
func (r *Redis) Close() error {
	return r.conn.Close()
}
