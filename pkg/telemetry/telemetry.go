/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package telemetry

import (
	"context"
	"log"
	"math/rand"
	"runtime"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

var (
	RequestCounter  metric.Int64Counter
	RequestLatency  metric.Float64Histogram
	RandomGauge     metric.Float64ObservableGauge
	MemoryHeapGauge metric.Int64ObservableGauge
	meterName       = "example.com/otel-prometheus"
)

func Init(ctx context.Context) {
	// OTLP Exporter
	otlpExp, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("failed to create OTLP exporter: %v", err)
	}

	// Periodic Reader
	reader := sdkmetric.NewPeriodicReader(otlpExp,
		sdkmetric.WithInterval(5*time.Second),
	)

	// Meter Provider
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
	)
	otel.SetMeterProvider(mp)

	// Meter
	meter := mp.Meter(meterName)

	// Counter
	RequestCounter, err = meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total HTTP requests received"),
	)
	if err != nil {
		log.Fatalf("failed to create counter: %v", err)
	}

	// Histogram
	RequestLatency, err = meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("HTTP request latency in seconds"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5),
	)
	if err != nil {
		log.Fatalf("failed to create histogram: %v", err)
	}

	// Gauge
	RandomGauge, err = meter.Float64ObservableGauge(
		"random_value",
		metric.WithDescription("A random gauge value"),
	)
	if err != nil {
		log.Fatalf("failed to create gauge: %v", err)
	}

	MemoryHeapGauge, err = meter.Int64ObservableGauge(
		"memory.heap",
		metric.WithDescription("Penggunaan memori heap saat ini"),
		metric.WithUnit("By"),
	)
	if err != nil {
		log.Fatalf("failed to create memory heap gauge: %v", err)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Callback untuk RandomGauge
	_, err = meter.RegisterCallback(func(ctx context.Context, obs metric.Observer) error {
		obs.ObserveFloat64(
			RandomGauge,
			rng.Float64()*100,
			metric.WithAttributes(attribute.String("unit", "random")),
		)
		return nil
	}, RandomGauge)
	if err != nil {
		log.Fatalf("failed to register gauge callback: %v", err)
	}

	// Callback untuk MemoryHeapGauge
	_, err = meter.RegisterCallback(func(ctx context.Context, obs metric.Observer) error {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		obs.ObserveInt64(
			MemoryHeapGauge,
			int64(m.HeapAlloc),
			metric.WithAttributes(attribute.String("unit", "bytes")),
		)
		return nil
	}, MemoryHeapGauge)
	if err != nil {
		log.Fatalf("failed to register memory heap gauge callback: %v", err)
	}
}
