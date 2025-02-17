package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/ashabykov/geospatial_cache_for_meetup/location"
)

type nearBy struct {
	Location location.Location `json:"location"`
	Radius   float64           `json:"radius"`
	Limit    int               `json:"limit"`
}

type tcase struct {
	name  string
	query nearBy
}

func tests() []tcase {
	inp := []float64{5000, 8000, 10000}
	ret := make([]tcase, 0, len(inp))
	for i := range inp {
		ret = append(ret, tcase{
			name: fmt.Sprintf("%f", inp[i]),
			query: nearBy{
				Location: location.Location{
					Name: "target",
					Lat:  43.244555,
					Lon:  76.940012,
					Ts:   location.Timestamp(time.Now().UTC().Unix()),
					TTL:  10 * time.Minute,
				},
				Radius: inp[i],
				Limit:  300,
			},
		})
	}
	return ret
}

func BenchmarkNearbyFunOutWrite(b *testing.B) {

	const url = "http://localhost:3000/nearby/v2"

	for _, tc := range tests() {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()

			for i := 0; i < b.N; i++ {

				c := http.Client{Timeout: time.Duration(3) * time.Second}

				data, err := json.Marshal(tc.query)
				if err != nil {
					b.Fatal(err)
				}

				resp, err := c.Post(url, "application/json", bytes.NewBuffer(data))
				if err != nil {
					b.Fatal(err)
					return
				}

				if resp.StatusCode != http.StatusOK {

					b.Fatalf("unexpected status code: %d", resp.StatusCode)
				}

				resp.Body.Close()
			}
		})
	}
}

func BenchmarkNearbyFunOutRead(b *testing.B) {

	const url = "http://localhost:3000/nearby/v1"

	for _, tc := range tests() {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()

			for i := 0; i < b.N; i++ {

				c := http.Client{Timeout: time.Duration(3) * time.Second}

				data, err := json.Marshal(tc.query)
				if err != nil {
					b.Fatal(err)
				}

				resp, err := c.Post(url, "application/json", bytes.NewBuffer(data))
				if err != nil {
					b.Fatal(err)
					return
				}

				if resp.StatusCode != http.StatusOK {

					b.Fatalf("unexpected status code: %d", resp.StatusCode)
				}

				resp.Body.Close()
			}
		})
	}
}
