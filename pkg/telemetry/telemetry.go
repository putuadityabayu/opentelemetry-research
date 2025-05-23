/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package telemetry

import (
	"context"
	"go.opentelemetry.io/otel/codes"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otlpmetric "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	otlptrace "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer trace.Tracer
	meter  metric.Meter
)

// InitTelemetry initializes OpenTelemetry with trace and metric exporters
func InitTelemetry(serviceName string, otelCollectorURL string) func() {
	ctx := context.Background()

	// Create a resource detailing service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("1.0.0"),
		),
	)
	if err != nil {
		log.Fatalf("Failed to create resource: %v", err)
	}

	// ===== TRACES =====
	// Create and configure trace exporter
	traceExporter, err := otlptrace.New(ctx,
		otlptrace.WithEndpoint(otelCollectorURL), // Pastikan port ini adalah port yang di-publish
		otlptrace.WithInsecure(),
		otlptrace.WithTimeout(30*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to create trace exporter: %v", err)
	}

	// Create trace provider with exporter
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(traceExporter),
	)
	otel.SetTracerProvider(tracerProvider)

	// Get a tracer
	tracer = tracerProvider.Tracer(serviceName)

	// ===== METRICS =====
	// Create and configure metric exporter
	metricExporter, err := otlpmetric.New(ctx,
		otlpmetric.WithEndpoint(otelCollectorURL), // Pastikan port ini adalah port yang di-publish
		otlpmetric.WithInsecure(),
		otlpmetric.WithTimeout(30*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to create metric exporter: %v", err)
	}

	// Create meter provider with exporter
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter, sdkmetric.WithInterval(1*time.Second))),
	)
	otel.SetMeterProvider(meterProvider)

	// Get a meter
	meter = meterProvider.Meter(serviceName)

	// Return cleanup function
	return func() {
		// Cleanup resources
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
		if err := meterProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down meter provider: %v", err)
		}
	}
}

// StartSpan starts a new span with the given name and returns the span with context
func StartSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return tracer.Start(ctx, spanName)
}

// AddSpanEvent adds an event to the current span
func AddSpanEvent(span trace.Span, name string, attributes map[string]string) {
	attrs := []attribute.KeyValue{}
	for k, v := range attributes {
		attrs = append(attrs, attribute.String(k, v))
	}
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// RecordSpanError records an error to the current span
func RecordSpanError(span trace.Span, err error, message string) {
	span.RecordError(err)
	span.SetStatus(codes.Error, message)
}

// CreateCounter creates and returns a new counter
func CreateCounter(name, description string, unit string) (metric.Int64Counter, error) {
	return meter.Int64Counter(
		name,
		metric.WithDescription(description),
		metric.WithUnit(unit),
	)
}

// CreateHistogram creates and returns a new histogram
func CreateHistogram(name, description string, unit string) (metric.Float64Histogram, error) {
	return meter.Float64Histogram(
		name,
		metric.WithDescription(description),
		metric.WithUnit(unit),
	)
}
