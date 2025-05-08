/*
 * Copyright (c) Portalnesia - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package main

import (
	"context"
	"errors"
	"github.com/subosito/gotenv"
	"log"
	"opentelemetry-research/pkg/telemetry"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdkTracer "go.opentelemetry.io/otel/sdk/trace"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
)

func main() {
	_ = gotenv.Load()

	ctx := context.Background()

	// Inisialisasi exporter untuk Google Cloud Trace
	exporter, err := texporter.New(texporter.WithProjectID("xentra-dummy-trial-01"))
	if err != nil {
		log.Fatalf("Gagal membuat exporter: %v", err)
	}

	// Inisialisasi TracerProvider dengan exporter
	tp := sdkTracer.NewTracerProvider(
		sdkTracer.WithBatcher(exporter),
	)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("Gagal shutdown TracerProvider: %v", err)
		}
	}()

	// Set global TracerProvider
	otel.SetTracerProvider(tp)

	// Buat tracer
	tracer := otel.Tracer("xentra-dummy")

	// Mulai span
	ctx, span := tracer.Start(ctx, "span-span-contoh")
	defer span.End()

	// Tambahkan atribut ke span
	span.SetAttributes(
		attribute.String("example.attribute", "value"),
		attribute.String("example.attribute2", "value"),
		attribute.String("example.attribute3", "value"),
	)

	// Simulasi pekerjaan
	time.Sleep(2 * time.Second)

	err = errors.New("contoh error")
	telemetry.RecordSpanError(span, err, "Faileddddd")

	log.Println("Tracing selesai.")
}
