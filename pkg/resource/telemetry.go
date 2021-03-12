package resource

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	exportmetric "go.opentelemetry.io/otel/sdk/export/metric"
	exporttrace "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"log"
)

var Int64ValueRecorder metric.Int64ValueRecorder
var Tracer trace.Tracer

func initTelemetry() {
	ctx := context.Background()

	driver := otlpgrpc.NewDriver(
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithEndpoint(OtlpEndpoint),
		//otlpgrpc.WithDialOption(grpc.WithBlock()), // useful for testing
	)

	otlpExporter, err := otlp.NewExporter(ctx, driver)
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	rs, err := resource.Detect(context.Background(), &resource.FromEnv{})
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	initTracer(otlpExporter, rs)
	initMetrics(otlpExporter, rs)
}

func initTracer(exporter exporttrace.SpanExporter, rs *resource.Resource) {
	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	tp := sdktrace.NewTracerProvider(sdktrace.WithConfig(
		sdktrace.Config{
			DefaultSampler: sdktrace.AlwaysSample(),
			Resource:       rs,
		}),
		sdktrace.WithSyncer(exporter))

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	Tracer = otel.GetTracerProvider().Tracer("io.treactor.tracing.golang", trace.WithInstrumentationVersion("0.5"))
}

func initMetrics(_ exportmetric.Exporter, _ *resource.Resource) {
	// TODO
}
