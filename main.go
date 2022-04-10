package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func InitTracer() func() {
	ctx := context.Background()

	otelAgentAddr := "0.0.0.0:4317"

	traceClient, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otelAgentAddr),
		otlptracegrpc.WithDialOption(),
	)
	if err != nil {
		log.Println(err)
	}

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String("Demo"),
		),
	)
	if err != nil {
		log.Println(err, "failed to create resource")
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceClient)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tracerProvider)

	return func() {
		cxt, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := traceClient.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}
}

func main() {
	shutDown := InitTracer()
	defer shutDown()

	newCtx, span := otel.Tracer("demo-tracer").Start(context.Background(), "Main function")

	fmt.Println("calling AddValue...")
	AddValue(newCtx)
	fmt.Println("calling AddValue2...")
	AddValue2(newCtx)
	fmt.Println("span end...")
	span.End()

	time.Sleep(time.Second * 10)
}

func AddValue(ctx context.Context) {
	_, span := otel.Tracer("demo-tracer").Start(ctx, "AddValue")
	defer span.End()
	span.SetAttributes(attribute.String("key1", "value"))
}

func AddValue2(ctx context.Context) {
	_, span := otel.Tracer("demo-tracer").Start(ctx, "AddValue2")
	defer span.End()
	span.SetAttributes(attribute.String("key2", "value2"))
}
