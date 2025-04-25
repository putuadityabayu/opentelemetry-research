/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

func main() {
	var n int
	var url string
	var delay time.Duration

	flag.IntVar(&n, "n", 10, "number of requests to send")
	flag.StringVar(&url, "url", "http://localhost:8080/", "target URL to hit")
	flag.DurationVar(&delay, "delay", 100*time.Millisecond, "delay between requests")
	flag.Parse()

	client := &http.Client{Timeout: 5 * time.Second}

	for i := 1; i <= n; i++ {
		start := time.Now()
		resp, err := client.Get(url)
		duration := time.Since(start)
		if err != nil {
			log.Printf("[%d] error: %v", i, err)
		} else {
			log.Printf("[%d] status: %s, duration: %v", i, resp.Status, duration)
			resp.Body.Close()
		}
		time.Sleep(delay)
	}
}

// go run client.go -n 100 -url http://localhost:8080/ -delay 250ms
