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

type Redis struct {
	conn *redis.Client
}

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
