package resource

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	exportmetric "go.opentelemetry.io/otel/sdk/export/metric"
	exporttrace "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"log"
	//sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
)

var Int64ValueRecorder metric.Int64ValueRecorder
var Tracer trace.Tracer

func initTelemetry() {
	stdoutExporter, err := stdout.NewExporter(stdout.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}
	otlpExporter, err := otlp.NewExporter(context.Background(), otlp.WithInsecure(), otlp.WithAddress(OtlpEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	_ = stdoutExporter

	rs, err := resource.Detect(context.Background(), &resource.FromEnv{})
	if err != nil {
		log.Fatal(err)
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
	Tracer = otel.Tracer("")
}

func initMetrics(exporter exportmetric.Exporter, rs *resource.Resource) {
	//
	//
	//sdkmetric.NewAccumulator()
	//// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	//// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	//tp := sdktrace.NewTracerProvider(sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	//	sdktrace.WithSyncer(exporter))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//global.SetTracerProvider(tp)
	//

	pusher := push.New(
		basic.New(
			simple.NewWithExactDistribution(),
			exporter,
		),
		exporter,
	)
	pusher.Start()
	defer pusher.Stop()
	// TODO: global.SetMeterProvider(pusher.MeterProvider())
}

func initMetricRecorders() {
	//meter := global.Meter("")
	//Int64ValueRecorder, _ = meter.NewInt64ValueRecorder("metrics")

}
