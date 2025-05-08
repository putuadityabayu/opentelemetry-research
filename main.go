/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package main

import (
	"context"
	"fmt"
	"github.com/subosito/gotenv"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"opentelemetry-research/pkg/telemetry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	requestCounter  metric.Int64Counter
	requestDuration metric.Float64Histogram
	errorCounter    metric.Int64Counter
)

func setupMetrics() error {
	var err error

	// Create a counter for total requests
	requestCounter, err = telemetry.CreateCounter(
		"api.request.total",
		"Total number of requests received",
		"{request}",
	)
	if err != nil {
		return fmt.Errorf("failed to create request counter: %w", err)
	}

	// Create a histogram for request duration
	requestDuration, err = telemetry.CreateHistogram(
		"api.request.duration",
		"Duration of requests in milliseconds",
		"ms",
	)
	if err != nil {
		return fmt.Errorf("failed to create duration histogram: %w", err)
	}

	// Create a counter for errors
	errorCounter, err = telemetry.CreateCounter(
		"api.error.total",
		"Total number of errors encountered",
		"{error}",
	)
	if err != nil {
		return fmt.Errorf("failed to create error counter: %w", err)
	}

	return nil
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	/*startTime := time.Now()

	// Record the request
	requestCounter.Add(ctx, 1, metric.WithAttributes(
		attribute.String("path", r.URL.Path),
		attribute.String("method", r.Method),
	))*/
	// Create a span for this request
	path := r.URL.Path
	if path == "/" {
		path = "index"
	}
	ctx, span := telemetry.StartSpan(ctx, path)
	defer span.End()

	// Add span attributes
	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.url", r.URL.String()),
	)

	// Simulate randomized response (success or error)
	simulateError := rand.Intn(100) < 80 // 30% chance of error

	if simulateError {
		// Simulate error handling
		err := CustomError{
			Data: "Contoh data",
			Name: "Ini Nama",
			Code: 553240,
		}
		span.AddEventHelper("error_data", map[string]any{
			"name": err.Name,
			"code": err.Code,
			"data": err.Data,
		})
		span.RecordErrorHelper(err, "Request processing failed")
		traceId, spanId := span.GetID()
		fmt.Printf("TraceID: %s, SpanID: %s\n", traceId, spanId)
		// Record error metric
		/*errorCounter.Add(ctx, 1, metric.WithAttributes(
			attribute.String("error_type", "internal_server_error"),
			attribute.String("path", r.URL.Path),
		))*/

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	} else {
		// Simulate successful response
		/*telemetry.AddSpanEvent(span, "request_processed", map[string]string{
			"status": "success",
		})*/

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Request processed successfully"))
	}

	// Record request duration
	//duration := float64(time.Since(startTime).Milliseconds())
	/*requestDuration.Record(ctx, duration, metric.WithAttributes(
		attribute.String("path", r.URL.Path),
		attribute.String("status", fmt.Sprintf("%d", w.Header().Get("Status"))),
	))*/
}

func main() {
	_ = gotenv.Load()
	// Initialize telemetry - using localhost for collector URL
	//cleanup := telemetry.InitTelemetry("demo-api-service", "localhost:4317")
	cleanup := telemetry.GcpTelemetry("demo-api-service")
	defer cleanup()

	// Setup metrics
	/*if err := setupMetrics(); err != nil {
		log.Fatalf("Failed to setup metrics: %v", err)
	}*/

	// Create HTTP handler with OpenTelemetry instrumentation
	handler := otelhttp.NewHandler(http.HandlerFunc(handleRequest), "server")
	http.Handle("/", handler)
	http.Handle("/user", handler)
	http.Handle("/product", handler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr: ":" + port,
	}

	// Handle graceful shutdown
	go func() {
		log.Printf("Server listening on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server gracefully stopped")
}

type CustomError struct {
	Data string
	Code int
	Name string
}

func (c CustomError) Error() string {
	return fmt.Sprintf("Name: %s, Code: %d, Data: %v", c.Name, c.Code, c.Data)
}

func (c CustomError) ToMap() map[string]string {
	return map[string]string{
		"name": c.Name,
		"code": fmt.Sprintf("%d", c.Code),
		"data": fmt.Sprintf("%v", c.Data),
	}
}
