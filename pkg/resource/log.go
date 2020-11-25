package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"os"
	"time"
)

type RLogger interface {
	InfoF(ctx context.Context, format string, a ...interface{})
	Info(ctx context.Context, message string)
	WarningF(ctx context.Context, format string, a ...interface{})
	Warning(ctx context.Context, message string)
	Error(ctx context.Context, r *http.Request, message string) string
	ErrorErr(ctx context.Context, r *http.Request, message string, err error) string
	ErrorF(ctx context.Context, r *http.Request, format string, a ...interface{}) string
	Flush()
}

type rPayLoad struct {
	EventTime      string          `json:"eventTime,omitempty"`
	ServiceContext rServiceContext `json:"serviceContext,omitempty"`
	Message        string          `json:"message,omitempty"`
	Context        rContext        `json:"context,omitempty"`
}

type rServiceContext struct {
	Service string `json:"service,omitempty"`
	Version string `json:"version,omitempty"`
}

type rContext struct {
	HttpRequest    rHttpRequest    `json:"httpRequest,omitempty"`
	User           string          `json:"user,omitempty"`
	ReportLocation rReportLocation `json:"reportLocation,omitempty"`
}

type rHttpRequest struct {
	Method             string `json:"method,omitempty"`
	Url                string `json:"url,omitempty"`
	UserAgent          string `json:"userAgent,omitempty"`
	Referrer           string `json:"referrer,omitempty"`
	ResponseStatusCode int    `json:"responseStatusCode,omitempty"`
	RemoteIp           string `json:"remoteIp,omitempty"`
}

type rReportLocation struct {
	FilePath     string `json:"filePath,omitempty"`
	LineNumber   int    `json:"lineNumber,omitempty"`
	FunctionName string `json:"functionName,omitempty"`
}


// https://cloud.google.com/logging/docs/agent/configuration#special-fields
type SLabel struct {
	LoggerName string `json:"loggerName,omitempty"`
}

type STimestamp struct {
	Seconds int64 `json:"seconds"`
	Nanos   int   `json:"nanos"`
}

type SServiceContext struct {
	Service string `json:"service"`
	Version string `json:"version"`
}

type SEntry struct {
	Severity       string           `json:"severity,omitempty"`
	Message        string           `json:"message,omitempty"`
	Timestamp      *STimestamp      `json:"timestamp,omitempty"`
	Labels         *SLabel          `json:"logging.googleapis.com/labels,omitempty"`
	ServiceContext *SServiceContext `json:"serviceContext,omitempty"`

	SpanId       string `json:"logging.googleapis.com/spanId,omitempty"`
	Trace        string `json:"logging.googleapis.com/trace,omitempty"`
	TraceSampled bool   `json:"logging.googleapis.com/trace_sampled,omitempty"`
}

type SLogger struct {
	projectId string
}

func NewSLogger(projectId string) *SLogger {
	return &SLogger{
		projectId: projectId,
	}
}

func (l *SLogger) InfoF(ctx context.Context, format string, a ...interface{}) {
	l.log(ctx, "INFO", fmt.Sprintf(format, a...))
}

func (l *SLogger) Info(ctx context.Context, message string) {
	l.log(ctx, "INFO", message)
}

func (l *SLogger) WarningF(ctx context.Context, format string, a ...interface{}) {
	l.log(ctx, "WARNING", fmt.Sprintf(format, a...))
}

func (l *SLogger) Warning(ctx context.Context, message string) {
	l.log(ctx, "WARNING", message)
}

func (l *SLogger) Error(ctx context.Context, r *http.Request, message string) string {
	l.log(ctx, "ERROR", message)
	return ""
}

func (l *SLogger) ErrorErr(ctx context.Context, r *http.Request, message string, err error) string {
	l.log(ctx, "ERROR", fmt.Sprintf("%s: %s", message, err.Error()))
	return ""
}

func (l *SLogger) ErrorF(ctx context.Context, r *http.Request, format string, a ...interface{}) string {
	l.log(ctx, "ERROR", fmt.Sprintf(format, a...))
	return ""
}

func (l *SLogger) Flush() {
	panic("implement me")
}

func (l *SLogger) log(ctx context.Context, severity string, message string) {
	now := time.Now()
	entry := &SEntry{
		Message:  message,
		Severity: severity,
		Timestamp: &STimestamp{
			Seconds: now.Unix(),
			Nanos:   now.Nanosecond(),
		},
		Labels: &SLabel{
			LoggerName: "treactor",
		},
		ServiceContext: &SServiceContext{
			Service: AppName,
			Version: AppVersion,
		},
	}
	entry = l.addSpan(ctx, entry)
	b, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
	os.Stdout.WriteString("\n")
}

func (l *SLogger) addSpan(ctx context.Context, entry *SEntry) *SEntry {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		entry.Trace = fmt.Sprintf("projects/%s/traces/%s", l.projectId, span.SpanContext().TraceID.String())
		entry.SpanId = span.SpanContext().SpanID.String()
	}
	return entry
}
