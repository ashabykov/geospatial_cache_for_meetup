package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/ashabykov/geospatial_cache_for_meetup/location"
)

type tcase struct {
	name  string
	query NearBy
}

func tests() []tcase {
	return []tcase{
		{
			name: "1000",
			query: NearBy{
				Location: location.Location{
					Name: "target",
					Lat:  43.244555,
					Lon:  76.940012,
					Ts:   location.Timestamp(time.Now().UTC().Unix()),
					TTL:  10 * time.Minute,
				},
				Radius: 1000,
				Limit:  100,
			},
		},
		{
			name: "5000",
			query: NearBy{
				Location: location.Location{
					Name: "target",
					Lat:  43.244555,
					Lon:  76.940012,
					Ts:   location.Timestamp(time.Now().UTC().Unix()),
					TTL:  10 * time.Minute,
				},
				Radius: 5000,
				Limit:  100,
			},
		},
		{
			name: "10000",
			query: NearBy{
				Location: location.Location{
					Name: "target",
					Lat:  43.244555,
					Lon:  76.940012,
					Ts:   location.Timestamp(time.Now().UTC().Unix()),
					TTL:  10 * time.Minute,
				},
				Radius: 10000,
				Limit:  100,
			},
		},
	}
}

func BenchmarkNearbyV2(b *testing.B) {

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

func BenchmarkNearbyV1(b *testing.B) {

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
