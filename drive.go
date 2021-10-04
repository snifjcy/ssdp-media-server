package main

import (
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
	"io/ioutil"
	"log"
	"math/rand"
	"media-server/google_drive"
	"media-server/server"
	"net/http"
)

func driveDecorator(driveSession *drive.Service, fileID string) func(func(http.ResponseWriter, *http.Request, int64, func(startOffset, endOffset int64) []byte)) http.HandlerFunc {
	file, err := google_drive.GetFile(driveSession, fileID)
	if err != nil {
		panic(err)
	}

	// Return a function that takes the server's videoHandler and returns a function that takes the server's actual http request
	return func(handler func(http.ResponseWriter, *http.Request, int64, func(int64, int64) []byte)) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, server.IDKey, rand.Int31())
			log.Printf("[%d] new request for file: %s, filesize: %d\n", ctx.Value(server.IDKey), file.Name, file.Size)
			getChunk := func(startOffset, endOffset int64) []byte {
				log.Printf("[%d] got chunk %10d-%10d\n", ctx.Value(server.IDKey), startOffset, endOffset)
				res, err := google_drive.DownloadFile(ctx, driveSession, fileID, startOffset, endOffset)
				if err != nil {
					if errors.Is(err, context.Canceled) {
						log.Printf("[%d] context canceled during getChunk\n", ctx.Value(server.IDKey))
						return []byte{}
					}
					log.Panicf("%+v\n", err)
				}
				chunkData, _ := ioutil.ReadAll(res.Body)
				return chunkData
			}
			handler(w, r.WithContext(ctx), file.Size, getChunk)
			log.Printf("[%d] end\n", ctx.Value(server.IDKey))
		}
	}
}
