package main

import (
	"context"
	"crypto/md5"
	"io"
	"net/http"
	"strconv"
	"time"
)

type CtxKey string

func DecryptBodyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" && r.Method != "PUT" {
			return
		}

		{
			timestamp, err := strconv.ParseInt(r.Header.Get("X-Timestamp"), 10, 64)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			offset := time.Now().UTC().Unix() - timestamp
			if offset < 0 || offset > 20 {
				w.WriteHeader(403)
				return
			}
		}

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		iv := md5.Sum([]byte(r.Header.Get("X-Timestamp") + SecretKey))
		plainData, err := Decrypt(string(reqBody), iv)

		ctx := context.WithValue(r.Context(), CtxKey("body"), plainData)
		ctx = context.WithValue(ctx, CtxKey("iv"), iv)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
