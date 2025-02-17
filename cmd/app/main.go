package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ashabykov/geospatial_cache_for_meetup/cmd"
	"github.com/ashabykov/geospatial_cache_for_meetup/location"
)

type NearBy struct {
	Location location.Location `json:"location"`
	Radius   float64           `json:"radius"`
	Limit    int               `json:"limit"`
}

func main() {

	ctx := cmd.WithContext(context.Background())

	client1, client2 := Init(ctx)

	go client2.SubscribeOnUpdates(ctx)

	r := chi.NewRouter()
	r.Post("/nearby/v1", func(w http.ResponseWriter, r *http.Request) {
		query := NearBy{}
		err := json.NewDecoder(r.Body).Decode(&query)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		results, err := client1.Near(query.Location, query.Radius, query.Limit)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		jsonItem, err := json.Marshal(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonItem)
	})

	r.Post("/nearby/v2", func(w http.ResponseWriter, r *http.Request) {
		query := NearBy{}
		err := json.NewDecoder(r.Body).Decode(&query)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		results, err := client2.Near(query.Location, query.Radius, query.Limit)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		jsonItem, err := json.Marshal(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonItem)
	})

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
