package resource

import (
	"os"
)

var (
	Port       string
	AppVersion string
	AppName    string

	OtlpEndpoint string

	Mode             string
	debug            string
	profile          string
	Base             string
	MaxBond          int
	tracePropagation string
	logMethod        string
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
func Configure() {
	// General Settings
	Port = getEnv("PORT", "3330")
	AppName = getEnv("SERVICE_NAME", "treactor")
	AppVersion = getEnv("SERVICE_VERSION", "0.0")
	// Reactor Specific Settings
	Mode = getEnv("TREACTOR_MODE", "local")
	// Reactor Fixed Settings
	Base = "/treact"
	MaxBond = 5

	OtlpEndpoint = getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "")

	tracePropagation = getEnv("TREACTOR_TRACE_PROPAGATION", "w3c")
	logMethod = os.Getenv("TREACTOR_LOG_METHOD")
}

func IsLocalMode() bool {
	return "local" == Mode
}

func IsKubernetesMode() bool {
	return "k8s" == Mode
}

func NextBond() string {
	return "n"
}

func TracePropagation() {
	//switch traceInternal {
	//case "b3":
	//	return &b3.HTTPFormat{}
	//default:
	//	return &b3.HTTPFormat{}
	//}
	return
}
