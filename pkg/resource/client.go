package resource

import (
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

var HttpClient *http.Client

func clientInit() {
	transport := otelhttp.NewTransport(http.DefaultTransport)
	HttpClient = &http.Client{Transport: transport}
}
