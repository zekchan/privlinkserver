package keyGenerators

import "github.com/satori/go.uuid"

func UUIDGenerator() string {
	return uuid.NewV4().String()
}
