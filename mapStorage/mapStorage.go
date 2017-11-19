package mapStorage

import (
	"errors"
	"time"
)


/**
	Простая in-memory реализация key-value стораджа с протуханием
	TODO: потенциально пофиксить возможные конфликты многопоточностти (map не потокобезопасный)
	TODO: наприсать еще пару реализаций (Например для Redis или блокчейна ¯\_(ツ)_/¯)
 */
type MapStorage map[string]string

func (v *MapStorage) Set(key string, url string) error {
	(*v)[key] = url
	return nil
}

func (v *MapStorage) Get(key string) (string, error) {
	str := *v
	url, ok := str[key]

	if ok {
		delete(str, key)
		return url, nil
	}

	return "", errors.New("no such key")
}

func (v *MapStorage) SetTTL(key string, ttl time.Duration) error {
	time.AfterFunc(ttl, func() {
		delete(*v, key)
	})
	return nil
}

func CreateMapStorage() *MapStorage {
	storage := make(MapStorage)

	return &storage
}