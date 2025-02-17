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

	r := chi.NewRouter()

	client := Init(ctx)

	go client.SubscribeOnUpdates(ctx)

	r.Post("/nearby/", func(w http.ResponseWriter, r *http.Request) {
		query := NearBy{}
		err := json.NewDecoder(r.Body).Decode(&query)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		results := client.Near(query.Location, query.Radius, query.Limit)

		jsonItem, err := json.Marshal(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(jsonItem)
	})

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
