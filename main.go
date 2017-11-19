package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recover"
	"github.com/iris-contrib/middleware/cors"
	"github.com/zekchan/privlinkserver/mapStorage"
	"github.com/satori/go.uuid"
	"strconv"
	"github.com/kataras/iris/context"
)

const DEFAULT_TTL int64 = 60 * 60
const DEFAULT_URL = "https://google.com"

func getTTL(ttl string) int64 {
	if time, err := strconv.ParseInt(ttl, 10, 64); ttl != "" && err == nil {
		return time
	}
	return DEFAULT_TTL
}

type Storage interface {
	Set(key string, url string) error
	Get(key string) (url string, err error)
	SetTTL(key string, ttl int64) error
}

func createLink(storage Storage) context.Handler {
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

		key := uuid.NewV4().String()

		storage.Set(key, url)
		storage.SetTTL(key, getTTL(ctx.FormValue("ttl")))

		ctx.JSON(Response{
			Key: key,
		})
	}
}
func redirect(storage Storage) context.Handler {

	return func(ctx iris.Context) {
		url, err := storage.Get(ctx.Params().Get("key"))

		if err == nil {
			ctx.Redirect(url)
		} else {
			ctx.Redirect(DEFAULT_URL)
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
	var storage Storage = mapStorage.CreateMapStorage()

	app.Post("/create", createLink(storage))
	app.Get("/{key}", redirect(storage))
	app.Run(iris.Addr(":8080"))
}
