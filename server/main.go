package server

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
)

type contextKey int

// BS : 2MB
const (
	IDKey contextKey = 0
	BS               = 1024 * 1024 * 1
)

func videoHandler(w http.ResponseWriter, r *http.Request, filesize int64, getChunk func(startOffset, endOffset int64) []byte) {
	ctx := r.Context()
	reqID := ctx.Value(IDKey)

	byteRange := r.Header.Get("Range")
	if byteRange == "" {
		log.Printf("[%d] byte range not set\n", reqID)
		addHeaders(&w, r.ProtoMajor, r.ProtoMinor, 0, filesize, filesize)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{})
	} else {
		strings := regexp.
			MustCompile("bytes=([0-9]+)-([0-9]+)?").
			FindStringSubmatch(byteRange)
		log.Printf("[%d] byte range set, %s-%s\n", reqID, strings[1], strings[2])

		startOffset, _ := strconv.ParseInt(strings[1], 10, 64)
		if startOffset == filesize {
			addHeaders(&w, r.ProtoMajor, r.ProtoMinor, filesize, filesize, filesize)
			w.WriteHeader(http.StatusPartialContent)
			w.Write([]byte{})
		} else if startOffset == 0 {
			shouldClose := addHeaders(&w, r.ProtoMajor, r.ProtoMinor, startOffset, filesize, filesize)
			w.WriteHeader(http.StatusOK)

			for true {
				safeEndOffset := int64(math.Min(float64(startOffset+BS), float64(filesize)))
				chunkData := getChunk(startOffset, safeEndOffset)
				select {
				case <-ctx.Done():
					log.Printf("[%d] context canceled while sending response\n", reqID)
					return
				default:
					w.Write(chunkData)
					log.Printf("[%d] successfully sent chunk (200)- %d\n", reqID, startOffset)
					startOffset += BS
					// return
				}
				if startOffset > safeEndOffset || shouldClose {
					return
				}
			}
		} else {
			shouldClose := addHeaders(&w, r.ProtoMajor, r.ProtoMinor, startOffset, filesize, filesize)
			w.WriteHeader(http.StatusPartialContent)

			for true {
				safeEndOffset := int64(math.Min(float64(startOffset+BS), float64(filesize)))
				chunkData := getChunk(startOffset, safeEndOffset)
				select {
				case <-ctx.Done():
					log.Printf("[%d] context canceled while sending response\n", reqID)
					return
				default:
					w.Write(chunkData)
					log.Printf("[%d] successfully sent chunk (206)- %d\n", reqID, startOffset)
					startOffset += BS
					// return
				}
				if startOffset > safeEndOffset || shouldClose {
					return
				}
			}
		}
	}

}

func addHeaders(w *http.ResponseWriter, ProtoMajor, ProtoMinor int, startOffset, endOffset, filesize int64) bool {
	(*w).Header().Add("Accept-Ranges", "bytes")
	(*w).Header().Add("Content-Disposition", "inline; \"test.mkv\"")
	(*w).Header().Add("Content-Length", fmt.Sprintf("%d", filesize))
	(*w).Header().Add("Content-Range", fmt.Sprintf("bytes %d-%d/%d", startOffset, endOffset, filesize))
	(*w).Header().Add("Content-Type", "video/mp4")

	if ProtoMajor == 1 && ProtoMinor == 1 {
		(*w).Header().Add("Connection", "close")
	}

	(*w).Header().Add("Access-Control-Allow-Origin", "*")
	(*w).Header().Add("Access-Control-Allow-Methods", "GET, OPTIONS")
	(*w).Header().Add("Access-Control-Allow-Headers", "Content-Type")

	(*w).Header().Add("transferMode.dlna.org", "Streaming")
	(*w).Header().Add("TimeSeekRange.dlna.org", "npt=0.00-")
	(*w).Header().Add("contentFeatures.dlna.org", "DLNA.ORG_OP=01;DLNA.ORG_CI=0;")

	return ProtoMajor == 1 && ProtoMinor == 1
}

// StartListening starts listening
func StartListening(handlerDec func(func(http.ResponseWriter, *http.Request, int64, func(startOffset, endOffset int64) []byte)) http.HandlerFunc) {
	http.HandleFunc("/video", handlerDec(videoHandler))

	const HOST = "192.168.1.95"
	const PORT = 8091
	log.Printf("Started listening on http://%s:%d", HOST, PORT)
	log.Panicln(http.ListenAndServe(fmt.Sprintf("%s:%d", HOST, PORT), nil))
}
