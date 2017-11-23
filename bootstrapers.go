package main

import (
	"github.com/zekchan/privlinkserver/mapStorage"
	"os"
	"github.com/satori/go.uuid"
	"log"
)

func getStorage() Storage {
	storages := map[string]func() Storage{
		"MAP": func() Storage {
			return mapStorage.CreateMapStorage()
		},
	}
	configStorage := os.Getenv("STORAGE")
	useStorage := "MAP"

	for key := range storages {
		if key == configStorage {
			useStorage = key
			break
		}
	}

	log.Printf("Used %v storage", useStorage)
	return storages[useStorage]()
}
func getKeyGenerator() func() string {
	generators := map[string]func() string{
		"UUID": func() string {
			return uuid.NewV4().String()
		},
	}
	configGenerator := os.Getenv("GENERATOR")
	useGenerator := "UUID"
	for key := range generators {
		if key == configGenerator {
			useGenerator = key
			break
		}
	}
	log.Printf("Used %v generator", useGenerator)
	return generators[useGenerator]

}
