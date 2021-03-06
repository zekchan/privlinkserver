package storages

import (
	"errors"
	"time"
)

type StoredData struct {
	value string;
	timer *time.Timer
}
type MapStorage map[string]StoredData

func (v *MapStorage) Set(key string, url string, ttl time.Duration) error {
	(*v)[key] = StoredData{
		value: url,
		timer: time.AfterFunc(ttl, func() {
			delete(*v, key)
		}),
	}
	return nil
}

func (v *MapStorage) Get(key string) (string, error) {
	str := *v
	data, ok := str[key]

	if ok {
		data.timer.Stop()
		delete(str, key)
		return data.value, nil
	}

	return "", errors.New("no such key")
}

func CreateMapStorage() *MapStorage {
	storage := make(MapStorage)

	return &storage
}
