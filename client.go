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
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

func getRandomPath() string {
	paths := []string{
		"/",
		"/user",
		"/product",
	}
	rand.Seed(time.Now().UnixNano())
	return paths[rand.Intn(len(paths))]
}

func main() {
	var n int
	var link string
	var delay time.Duration

	flag.IntVar(&n, "n", 20, "number of requests to send")
	flag.StringVar(&link, "url", "http://localhost:8080", "target URL to hit")
	flag.DurationVar(&delay, "delay", 100*time.Millisecond, "delay between requests")
	flag.Parse()

	client := &http.Client{Timeout: 5 * time.Second}

	for i := 1; i <= n; i++ {
		start := time.Now()
		u, _ := url.Parse(link)
		u.Path = getRandomPath()
		resp, err := client.Get(u.String())
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
