package main

import (
	"bytes"
	"context"
	"image"
	"image/draw"
	"time"

	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	redsyncredislib "github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

const defaultTimeout = 5 * time.Second

// Redis is a utility client for reading and writing values to the Redis server.
type Redis struct {
	Client     *redis.Client
	Pool       *redsyncredis.Pool
	SyncClient *redsync.Redsync
}

// Connect connects to the Redis server using the configuration values provided.
func (r *Redis) Connect(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)

	defer cancel()

	opts, err := redis.ParseURL(url)

	if err != nil {
		return err
	}

	r.Client = redis.NewClient(opts)

	if err := r.Client.Ping(ctx).Err(); err != nil {
		return err
	}

	pool := redsyncredislib.NewPool(r.Client)

	r.SyncClient = redsync.New(pool)

	return nil
}

// GetBytes gets the value from Redis by the key and returns the value as a byte array.
func (r *Redis) GetBytes(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	existsResult := r.Client.Exists(ctx, key)

	if err := existsResult.Err(); err != nil {
		return nil, err
	}

	if existsResult.Val() == 0 {
		return nil, nil
	}

	result := r.Client.Get(ctx, key)

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

	result := r.Client.Exists(ctx, key)

	return result.Val() == 1, result.Err()
}

// Delete deletes the key from the database if it exists.
func (r *Redis) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	return r.Client.Del(ctx, key).Err()
}

// Set puts the key-value into the database with an optional TTL value.
func (r *Redis) Set(key string, value interface{}, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	return r.Client.Set(ctx, key, value, ttl).Err()
}

// NewMutex creates a new mutually exclusive lock that only one process can hold.
func (r *Redis) NewMutex(name string) *Mutex {
	if r.Client == nil || r.SyncClient == nil {
		return &Mutex{
			Mutex: nil,
		}
	}

	return &Mutex{
		Mutex: r.SyncClient.NewMutex(name),
	}
}

// Close closes the connection to the database.
func (r *Redis) Close() error {
	return r.Client.Close()
}

// Mutex is a mutually exclusive lock held across all processes.
type Mutex struct {
	Mutex *redsync.Mutex
}

// Lock will lock the mutex so no other process can hold it.
func (m *Mutex) Lock() error {
	if m.Mutex == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	return m.Mutex.LockContext(ctx)
}

// Unlock will allow any other process to obtain a lock with the same key.
func (m *Mutex) Unlock() error {
	if m.Mutex == nil {
		return nil
	}

	_, err := m.Mutex.Unlock()

	return err
}
