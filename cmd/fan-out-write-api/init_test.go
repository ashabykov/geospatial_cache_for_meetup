package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func BenchmarkInit(b *testing.B) {

	c := http.Client{Timeout: time.Duration(1) * time.Second}

	data := bytes.NewBuffer([]byte(`{"location":{"lon":76.940012,"lat":43.244555},"radius":5000,"limit":100}`))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {

			resp, err := c.Post("http://localhost:3000/nearby", "application/json", data)
			if err != nil {
				fmt.Errorf("Error %s", err)
				return
			}

			_, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Errorf("Error %s", err)
			}
		}
	})
}
