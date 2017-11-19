package mapStorage

import "errors"

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

func (v *MapStorage) SetTTL(key string, ttl int64) error {
	// TODO: implement
	return nil
}

func CreateMapStorage() *MapStorage {
	storage := make(MapStorage)

	return &storage
}