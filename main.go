package main

import (
	"fmt"
	"log"
	"strconv"
	"os"
	"encoding/json"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

const dataDir = "./data/"
const tokensFile = dataDir + "tokens.bin"
const address = ":1264"
const suffix = ".s1"
var authPrefix = []byte("Bearer ")

func indexRoute(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("GET    /token\n")
	ctx.WriteString("GET    /db/\n")
	ctx.WriteString("GET    /db/:key\n")
	ctx.WriteString("PUT    /db/:key\n")
	ctx.WriteString("DELETE /db/:key")
}

func tokenRoute(ctx *fasthttp.RequestCtx) {
	token := generateToken()
	ctx.WriteString(token)
}

func dbRouteGet(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	token, ok := getToken(ctx)

	if (ok) {
		key := ctx.UserValue("key")
		strKey, ok := key.(string)
		if !ok {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusBadRequest), fasthttp.StatusBadRequest)
			return
		}

		data := getData(token, strKey)
		if data != nil {
			ctx.Write(data)
		} else {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
		}
		return
	}

	ctx.Response.Header.Set("WWW-Authenticate", "Bearer realm=Restricted")
	ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
}

func dbRoutePut(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	token, ok := getToken(ctx)

	if (ok) {
		key := ctx.UserValue("key")
		strKey, ok := key.(string)
		if !ok {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusBadRequest), fasthttp.StatusBadRequest)
			return
		}

		storeData(token, strKey, ctx.PostBody())
		ctx.SetStatusCode(fasthttp.StatusNoContent)
		return
	}

	ctx.Response.Header.Set("WWW-Authenticate", "Bearer realm=Restricted")
	ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
}

func dbRouteDelete(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	token, ok := getToken(ctx)

	if (ok) {
		key := ctx.UserValue("key")
		strKey, ok := key.(string)
		if !ok {
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusBadRequest), fasthttp.StatusBadRequest)
			return
		}

		deleteData(token, strKey)
		ctx.SetStatusCode(fasthttp.StatusNoContent)
		return
	}

	ctx.Response.Header.Set("WWW-Authenticate", "Bearer realm=Restricted")
	ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
}

func dbGenericRouteGet(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	token, ok := getToken(ctx)

	if (ok) {
		keys := getKeys(token)
		if keys == nil {
			keys = []string{}
		}
		res, err := json.Marshal(keys)
		if err == nil {
			ctx.Write(res)
			return
		}

		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusInternalServerError), fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.Header.Set("WWW-Authenticate", "Bearer realm=Restricted")
	ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
}

func main() {
	r := router.New()

	r.GET("/", indexRoute)
	r.GET("/token", tokenRoute)

	r.GET("/db/", dbGenericRouteGet)

	r.GET("/db/{key}", dbRouteGet)
	r.PUT("/db/{key}", dbRoutePut)
	r.DELETE("/db/{key}", dbRouteDelete)

	fmt.Println("Starting server on " + address)
	log.Fatal(fasthttp.ListenAndServe(address, r.Handler))
}

func init() {
	os.MkdirAll(dataDir, os.ModePerm)

	fmt.Println("Loading tokens...")
	tokensFileExists := exists(tokensFile)
	if tokensFileExists {
		loadTokens()
	}
	fmt.Println("Loaded " + strconv.Itoa(len(tokens)) + " tokens!")

	fmt.Println("")
}
