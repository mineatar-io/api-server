package store

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

type FileStore struct {
	cancelCtx     context.Context
	cancelFunc    context.CancelFunc
	BaseDir       string
	CleanInterval time.Duration
}

func (s *FileStore) Initialize(config map[string]interface{}) error {
	baseDir, ok := config["dir"].(string)

	if !ok {
		return fmt.Errorf("filestore: invalid base directory value: %s", config["dir"])
	}

	if err := os.MkdirAll(path.Clean(baseDir), 0777); err != nil {
		return err
	}

	s.BaseDir = path.Clean(baseDir)

	cleanIntervalValue, ok := config["clean_interval"].(string)

	if !ok {
		return fmt.Errorf("filestore: invalid clean interval value: %s", config["clean_interval"])
	}

	cleanInterval, err := time.ParseDuration(cleanIntervalValue)

	if err != nil {
		return err
	}

	s.CleanInterval = cleanInterval
	s.cancelCtx, s.cancelFunc = context.WithCancel(context.Background())

	go s.cleanupGoroutine()

	return nil
}

func (s *FileStore) cleanupGoroutine() {
	if err := s.runCleanup(); err != nil {
		log.Println(err)
	}

loop:
	for {
		select {
		case <-s.cancelCtx.Done():
			break loop
		case <-time.After(s.CleanInterval):
			{
				if err := s.runCleanup(); err != nil {
					log.Println(err)
				}

				break
			}
		}
	}
}

func (s *FileStore) runCleanup() error {
	files, err := os.ReadDir(s.BaseDir)

	if err != nil {
		return err
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".expiration.txt") {
			continue
		}

		data, err := os.ReadFile(path.Join(s.BaseDir, file.Name()))

		if err != nil {
			log.Println(err)

			continue
		}

		expirationDate, err := time.Parse(time.RFC3339, string(data))

		if err != nil {
			log.Println(err)

			continue
		}

		if !time.Now().After(expirationDate) {
			continue
		}

		if err = os.Remove(path.Join(s.BaseDir, file.Name())); err != nil {
			log.Println(err)
		}

		if err = os.Remove(path.Join(s.BaseDir, strings.Replace(file.Name(), ".expiration.txt", ".bin", 1))); err != nil {
			log.Println(err)
		}
	}

	return nil
}

func (s *FileStore) GetBytes(key string) ([]byte, bool, error) {
	expiration, err := os.ReadFile(path.Join(s.BaseDir, fmt.Sprintf("%s.expiration.txt", key)))

	if err == nil {
		expirationDate, err := time.Parse(time.RFC3339, string(expiration))

		if err == nil && time.Now().After(expirationDate) {
			return nil, false, nil
		}
	}

	data, err := os.ReadFile(path.Join(s.BaseDir, fmt.Sprintf("%s.bin", key)))

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, false, nil
		}

		return nil, false, err
	}

	return data, true, nil
}

func (s *FileStore) GetNRGBA(key string) (*image.NRGBA, bool, error) {
	data, exists, err := s.GetBytes(key)

	if !exists || err != nil {
		return nil, exists, err
	}

	img, format, err := image.Decode(bytes.NewReader(data))

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

func (s *FileStore) Exists(key string) (bool, error) {
	expiration, err := os.ReadFile(path.Join(s.BaseDir, fmt.Sprintf("%s.expiration.txt", key)))

	if err == nil {
		expirationDate, err := time.Parse(time.RFC3339, string(expiration))

		if err == nil && time.Now().After(expirationDate) {
			return false, nil
		}
	}

	if _, err = os.Stat(path.Join(s.BaseDir, fmt.Sprintf("%s.bin", key))); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *FileStore) SetBytes(key string, data []byte, ttl time.Duration) error {
	if err := os.WriteFile(path.Join(s.BaseDir, fmt.Sprintf("%s.bin", key)), data, 0777); err != nil {
		return err
	}

	if ttl > 0 {
		if err := os.WriteFile(path.Join(s.BaseDir, fmt.Sprintf("%s.expiration.txt", key)), []byte(time.Now().Add(ttl).Format(time.RFC3339)), 0777); err != nil {
			return err
		}
	} else {
		if err := os.RemoveAll(path.Join(s.BaseDir, fmt.Sprintf("%s.expiration.txt", key))); err != nil {
			return err
		}
	}

	return nil
}

func (s *FileStore) Delete(key string) error {
	if err := os.RemoveAll(path.Join(s.BaseDir, fmt.Sprintf("%s.bin", key))); err != nil {
		return err
	}

	if err := os.RemoveAll(path.Join(s.BaseDir, fmt.Sprintf("%s.expiration.txt", key))); err != nil {
		return err
	}

	return nil
}

func (s *FileStore) Close() error {
	s.cancelFunc()

	return nil
}

var _ Store = &FileStore{}
