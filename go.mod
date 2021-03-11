module github.com/treactor/treactor-go

go 1.14

require (
	github.com/golang/protobuf v1.4.2
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.18.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.18.0
	go.opentelemetry.io/otel v0.18.0
	go.opentelemetry.io/otel/exporters/otlp v0.18.0
	go.opentelemetry.io/otel/exporters/stdout v0.18.0 // indirect
	go.opentelemetry.io/otel/metric v0.18.0
	go.opentelemetry.io/otel/sdk v0.18.0
	go.opentelemetry.io/otel/sdk/export/metric v0.18.0
	go.opentelemetry.io/otel/sdk/metric v0.18.0 // indirect
	go.opentelemetry.io/otel/trace v0.18.0
	golang.org/x/net v0.0.0-20201021035429-f5854403a974
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0 // indirect
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v2 v2.3.0
)
