package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
)

func main() {
	// fmt.Println(GetTrainSchedule("CNP", "KM", "20-01-2026"))
	// fmt.Println(KaiWebSchedule("CNP", "KM", "20-01-2026"))
	// fmt.Println(TravelokaTrainsSchedule("CNP", "KM", "20-01-2026"))
	r := chi.NewRouter()
	r.Use(httprate.Limit(
		1500,
		time.Minute*10,
		httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
			return r.Header.Get("X-Device-Id"), nil
		}),
	))

	{
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
		})
		r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(405)
		})
	}

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok!"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(DecryptBodyMiddleware)

		r.Route("/schedule", func(r chi.Router) {
			r.Post("/is-available", func(w http.ResponseWriter, r *http.Request) {
				body := r.Context().Value(CtxKey("body")).(string)

				var reqPayload = TrainScheduleReqBody{}
				if err := json.Unmarshal([]byte(body), &reqPayload); err != nil {
					w.WriteHeader(400)
					return
				}

				cacheKey := reqPayload.Org + reqPayload.Dest + reqPayload.DepartDate

				mu := GetMutexForKey(cacheKey)
				mu.Lock()
				defer mu.Unlock()

				var (
					trains     []Train
					cachedData any
					found      bool
					err        error
				)

				cachedData, found = c.Get(cacheKey)
				if found {
					trains = cachedData.([]Train)
				} else {
					trains, err = GetTrainSchedule(reqPayload.Org, reqPayload.Dest, reqPayload.DepartDate)
					if err != nil {
						log.Printf("Failed get train schedule: %v\n", err)
						w.WriteHeader(500)
						return
					}
					c.SetDefault(cacheKey, trains)
					fmt.Println("generated: " + cacheKey)
				}

				for _, t := range trains {
					// fmt.Println(reqPayload.TrainId, "-", strings.ToUpper(t.Name+t.Class), "-", t.Stock)
					if reqPayload.TrainId == strings.ToUpper(t.Name+t.Class) && t.Stock >= reqPayload.TotPsg {
						w.WriteHeader(200)
						return
					}
				}

				w.WriteHeader(404)
			})

			r.Put("/", func(w http.ResponseWriter, r *http.Request) {
				body := r.Context().Value(CtxKey("body")).(string)

				var reqPayload = UpdateTrainScheduleReqBody{}
				if err := json.Unmarshal([]byte(body), &reqPayload); err != nil {
					w.WriteHeader(400)
					return
				}

				cacheKey := reqPayload.Org + reqPayload.Dest + reqPayload.DepartDate

				mu := GetMutexForKey(cacheKey)
				mu.Lock()
				defer mu.Unlock()

				var (
					trains     []Train
					cachedData any
					found      bool
					expiredAt  int64
					err        error
				)

				cachedData, found = c.Get(cacheKey)
				if found {
					expiredAt = c.Items()[cacheKey].Expiration
					trains = cachedData.([]Train)
				} else {
					trains, err = GetTrainSchedule(reqPayload.Org, reqPayload.Dest, reqPayload.DepartDate)
					if err != nil {
						log.Printf("Failed get train schedule: %v\n", err)
						w.WriteHeader(500)
						return
					}
					expiredAt = time.Now().UnixNano() + ((time.Duration(1) * time.Minute).Nanoseconds())
					fmt.Println("generated: " + strconv.Itoa(int(expiredAt)))
				}

				found = false
				for n, t := range trains {
					if reqPayload.TrainId == t.Name+t.Class {
						trains[n].Stock = reqPayload.Stock
						found = true
						break
					}
				}

				if !found {
					trains = append(trains, Train{
						Name:  reqPayload.TrainId,
						Stock: reqPayload.Stock,
					})
				}

				c.Set(cacheKey, trains, time.Duration(expiredAt)*time.Nanosecond)
			})
		})
	})

	http.ListenAndServe(os.Args[1], r)
}
