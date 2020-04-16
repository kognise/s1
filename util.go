package main

import (
	"os"
	"bytes"
	
	"github.com/valyala/fasthttp"
)

func exists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func getToken(ctx *fasthttp.RequestCtx) (string, bool) {
	auth := ctx.Request.Header.Peek("Authorization")
	if bytes.HasPrefix(auth, authPrefix) {
		token := string(auth[len(authPrefix):])

		if _, ok := tokens[token]; ok {
			return token, true
		}
	}

	return "", false
}