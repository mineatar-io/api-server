package redis

import (
	"bytes"
	"context"
	"image"
	"image/draw"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	redisGetMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: "redis_get_count",
		Help: "The amount of Redis GET requests",
	})
	redisSetMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: "redis_set_count",
		Help: "The amount of Redis SET requests",
	})
	redisExistsMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: "redis_exists_count",
		Help: "The amount of Redis EXIST requests",
	})
	redisDeleteMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: "redis_delete_count",
		Help: "The amount of Redis DELETE requests",
	})
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
	redisGetMetric.Inc()

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
	redisGetMetric.Inc()

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

func (r *Redis) GetNRGBA(key string) (*image.NRGBA, bool, error) {
	redisGetMetric.Inc()

	value, ok, err := r.GetBytes(key)

	if err != nil {
		return nil, false, err
	}

	if !ok {
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
	redisExistsMetric.Inc()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	result := r.conn.Exists(ctx, key)

	return result.Val() == 1, result.Err()
}

func (r *Redis) Delete(key string) error {
	redisDeleteMetric.Inc()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	return r.conn.Del(ctx, key).Err()
}

func (r *Redis) Set(key string, value interface{}, ttl time.Duration) error {
	redisSetMetric.Inc()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	return r.conn.Set(ctx, key, value, ttl).Err()
}

func (r *Redis) Close() error {
	return r.conn.Close()
}
