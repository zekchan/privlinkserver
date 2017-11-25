package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recover"
	"github.com/iris-contrib/middleware/cors"
	"strconv"
	"github.com/kataras/iris/context"
	"github.com/zekchan/privlinkserver/utils"
	"time"
	"os"
	"log"
	"github.com/zekchan/privlinkserver/keyGenerators"
	"github.com/zekchan/privlinkserver/storages"
)

const DEFAULT_TTL time.Duration = time.Hour

func getTTL(ttl string) time.Duration {
	if parsedTime, err := strconv.ParseInt(ttl, 10, 64); ttl != "" && err == nil {
		return time.Duration(parsedTime) * time.Second
	}
	return DEFAULT_TTL
}

type Storage interface {
	Set(key string, url string, ttl time.Duration) error
	Get(key string) (url string, err error)
}

func getStorage() Storage {
	storages := map[string]func() Storage{
		"MAP": func() Storage {
			return storages.CreateMapStorage()
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
		"UUID": keyGenerators.UUIDGenerator,
		"RAND": keyGenerators.RandomKeyGenerator,
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
func createLink(storage Storage, keyGenerator func() string) context.Handler {
	type Response struct {
		Key string `json:"key"`
	}
	badResponse := struct {
		Error string `json:"error"`
	}{Error: "Bad request"}
	return func(ctx iris.Context) {
		url := ctx.FormValue("url")

		if url == "" {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(badResponse)
			return
		}

		key := keyGenerator()

		storage.Set(key, url, getTTL(ctx.FormValue("ttl")))

		ctx.JSON(Response{
			Key: key,
		})
	}
}
func redirect(storage Storage) context.Handler {
	defaultUrl := utils.EnvOr("DEFAULT_REDIRECT", "https://google.com")
	return func(ctx iris.Context) {
		url, err := storage.Get(ctx.Params().Get("key"))

		if err == nil {
			ctx.Redirect(url)
		} else {
			ctx.Redirect(defaultUrl)
		}
	}
}

func main() {
	app := iris.New()
	app.Use(recover.New())
	app.WrapRouter(cors.WrapNext(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	}))
	storage := getStorage()
	keyGenerator := getKeyGenerator()
	var port = utils.EnvOr("PORT", "80")

	app.Post("/create", createLink(storage, keyGenerator))
	app.Get("/{key}", redirect(storage))
	app.Run(iris.Addr(":" + port))
}
