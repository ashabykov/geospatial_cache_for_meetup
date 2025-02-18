package main

import (
	"context"
	"expvar"
	"log"
	"net/http"
	"net/http/pprof"

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

	h := New(ctx)

	h.Star(ctx)

	r := chi.NewRouter()

	r.Post("/nearby/v1", h.FanOutReadClientHandler)
	r.Post("/nearby/v2", h.FanOutWriteClientHandler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/pprof/", http.StatusMovedPermanently)
	})
	r.HandleFunc("/pprof", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/", http.StatusMovedPermanently)
	})

	r.HandleFunc("/pprof/*", pprof.Index)
	r.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/pprof/profile", pprof.Profile)
	r.HandleFunc("/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/pprof/trace", pprof.Trace)
	r.Handle("/vars", expvar.Handler())

	r.Handle("/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/pprof/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/pprof/mutex", pprof.Handler("mutex"))
	r.Handle("/pprof/heap", pprof.Handler("heap"))
	r.Handle("/pprof/block", pprof.Handler("block"))
	r.Handle("/pprof/allocs", pprof.Handler("allocs"))

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
