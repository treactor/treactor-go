package treact

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/treactor/treactor-go/pkg/chem"
	"github.com/treactor/treactor-go/pkg/execute"
	"github.com/treactor/treactor-go/pkg/resource"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/label"
	trace "go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	InsertId string
}

func executePlan(w http.ResponseWriter, r *http.Request, ctx context.Context, plan execute.Plan) {
	ch := make(chan execute.Capture, plan.Calls())
	plan.Execute(ctx, ch)

	elems := len(ch)
	capture := execute.Capture{
		Name:    resource.AppName,
		Headers: make(map[string]string, len(r.Header)),
		Bonds:   make([]execute.Capture, elems),
	}
	for key, values := range r.Header {
		capture.Headers[key] = strings.Join(values, "|")
	}
	for i := 0; i < elems; i++ {
		capture.Bonds[i] = <-ch
	}
	bytes, _ := json.MarshalIndent(capture, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func failure(ctx context.Context, w http.ResponseWriter, r *http.Request, message string, err error) {
	insertId := resource.Logger.ErrorErr(ctx, r, message, err)
	errorResponse := &ErrorResponse{
		InsertId: insertId,
	}
	w.WriteHeader(400)
	w.Header().Set("Content-Type", "application/json")
	bytes, _ := json.MarshalIndent(errorResponse, "", "\t")
	w.Write(bytes)
}

func TReactSplitHandle(w http.ResponseWriter, r *http.Request) {
	ctx, span := resource.Tracer.Start(r.Context(), "TReactSplitHandle", trace.WithAttributes(label.KeyValue{
		Key: "x",
		Value: label.Value{

		},
	}))
	defer span.End()
	//_, span := trace.StartSpan(r.Context(), "split.Get")
	//defer span.End()
	//span.Annotate([]trace.Attribute{trace.StringAttribute("key", "value")}, "something happened")
	//span.AddAttributes(trace.StringAttribute("hello", "world"))
	url := r.URL
	molecule := url.Query().Get("molecule")
	resource.Logger.InfoF(ctx, "Starting reaction for molecule %s", molecule)

	plan, err := execute.Parse(molecule)
	if err != nil {
		failure(ctx, w, r, "Unable to parse molecule", err)
		return
	}

	executePlan(w, r, ctx, plan)
	resource.Logger.WarningF(ctx, "Cooling down reaction, finished %s", molecule)
}

func TReactBondHandle(w http.ResponseWriter, r *http.Request) {
	ctx, span := resource.Tracer.Start(r.Context(), "TReactBondHandle")
	defer span.End()
	url := r.URL
	plan, err := execute.Parse(url.Query().Get("molecule"))
	if err != nil {
		failure(r.Context(), w, r, "Unable to parse molecule", err)
		return
	}
	executePlan(w, r, ctx, plan)
}

func TReactAtomHandle(w http.ResponseWriter, r *http.Request) {
	ctx, span := resource.Tracer.Start(r.Context(), "TReactAtomHandle")
	defer span.End()

	resource.Int64ValueRecorder.Measurement(12)

	url := r.URL
	symbol := url.Query().Get("symbol")
	block, err := execute.ParseBlock(symbol)
	if err != nil {
		failure(r.Context(), w, r, "Unable to parse atom", err)
		return
	}

	atom := resource.Atoms.Symbols[block.Block]

	var mb []byte
	if block.KV["mem"] != "" {
		mb = mem(ctx, block.KV["mem"])
	}

	if block.KV["cpu"] != "" {
		cpu(ctx, block.KV["cpu"])
	}

	span.AddEvent("AtomEvent",
		// TODO: label.Int("foo", 12)
	)

	resource.Logger.InfoF(r.Context(), "Atom %s (%d)", atom.Name, atom.Number)
	capture := execute.Capture{
		Name:    resource.AppName,
		Headers: make(map[string]string, len(r.Header)),
	}
	for key, values := range r.Header {
		capture.Headers[key] = strings.Join(values, "|")
	}
	bytes, _ := json.MarshalIndent(capture, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
	_ = mb
}

func ReactorHealthz(_ http.ResponseWriter, _ *http.Request) {
}

func Serve() {
	atoms := chem.NewAtoms()

	fmt.Printf("Telemetry Reactor (%s:%s) listening on port %s\n", resource.AppName, resource.AppVersion, resource.Port)
	fmt.Printf("Mode: %s\n", resource.Mode)

	r := http.NewServeMux()
	r.HandleFunc("/healthz", ReactorHealthz)
	r.HandleFunc(fmt.Sprintf("%s/split", resource.Base), TReactSplitHandle)
	r.HandleFunc(fmt.Sprintf("%s/bond/1", resource.Base), TReactBondHandle)
	r.HandleFunc(fmt.Sprintf("%s/bond/2", resource.Base), TReactBondHandle)
	r.HandleFunc(fmt.Sprintf("%s/bond/3", resource.Base), TReactBondHandle)
	r.HandleFunc(fmt.Sprintf("%s/bond/4", resource.Base), TReactBondHandle)
	r.HandleFunc(fmt.Sprintf("%s/bond/5", resource.Base), TReactBondHandle)
	r.HandleFunc(fmt.Sprintf("%s/bond/n", resource.Base), TReactBondHandle)
	for sym := range atoms.Symbols {
		r.HandleFunc(fmt.Sprintf("%s/atom/%s", resource.Base, strings.ToLower(sym)), TReactAtomHandle)
	}
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", resource.Port), otelhttp.NewHandler(
		r, "TReact")))
}
