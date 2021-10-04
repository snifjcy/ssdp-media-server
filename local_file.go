package main

import (
	"golang.org/x/net/context"
	"math/rand"
	"media-server/server"
	"net/http"
	"os"
)

func localFileDecorator(path string) func(func(http.ResponseWriter, *http.Request, int64, func(startOffset, endOffset int64) []byte)) http.HandlerFunc {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	stats, err := f.Stat()
	if err != nil {
		panic(err)
	}

	// Return a function that takes the server's videoHandler and returns a function that takes the server's actual http request
	return func(handler func(http.ResponseWriter, *http.Request, int64, func(int64, int64) []byte)) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, server.IDKey, rand.Int31())

			getChunk := func(startOffset, endOffset int64) []byte {
				chunkData := make([]byte, server.BS)
				f.Seek(startOffset, 0)
				f.Read(chunkData)
				return chunkData
			}
			handler(w, r.WithContext(ctx), stats.Size(), getChunk)
		}
	}
}
