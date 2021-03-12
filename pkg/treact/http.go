package treact

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/protobuf/encoding/protojson"

	treactorpb "github.com/treactor/treactor-go/io/treactor/v1alpha"
	"github.com/treactor/treactor-go/pkg/element"
	"github.com/treactor/treactor-go/pkg/execute"
	"github.com/treactor/treactor-go/pkg/resource"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	trace "go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	InsertId string
}

func executePlan(w http.ResponseWriter, r *http.Request, ctx context.Context, plan execute.Plan) {
	ch := make(chan *treactorpb.Bond, plan.Calls())
	plan.Execute(ctx, ch)

	elems := len(ch)

	node := &treactorpb.Node{
		Name:      resource.AppName,
		Version:   resource.AppVersion,
		Framework: resource.Framework,
		Request: &treactorpb.TReactorRequest{
			Path:    r.RequestURI,
			Headers: make(map[string]string, len(r.Header)),
		},
		Bonds: make([]*treactorpb.Bond, elems),
		Atom:  nil,
	}
	for key, values := range r.Header {
		node.Request.Headers[key] = strings.Join(values, "|")
	}
	for i := 0; i < elems; i++ {
		node.Bonds[i] = <-ch
	}
	bytes, _ := protojson.Marshal(node)
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
	ctx, span := resource.Tracer.Start(r.Context(), "TReactSplitHandle", trace.WithAttributes(
		attribute.String("x", "foo")))
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

	atom := resource.Atoms.ElementByName[strings.ToLower(block.Block)]

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

	node := &treactorpb.Node{
		Name:      resource.AppName,
		Version:   resource.AppVersion,
		Framework: resource.Framework,
		Request: &treactorpb.TReactorRequest{
			Path:    r.RequestURI,
			Headers: make(map[string]string, len(r.Header)),
		},
		Bonds: nil,
		Atom: &treactorpb.Atom{
			Number: resource.Number,
			Symbol: atom.Symbol,
			Name:   atom.Name,
			Period: &atom.Period,
			Group:  &atom.Group,
		},
	}
	for key, values := range r.Header {
		node.Request.Headers[key] = strings.Join(values, "|")
	}

	bytes, _ := protojson.Marshal(node)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
	_ = mb
}

func TReactInfoHandle(w http.ResponseWriter, r *http.Request) {
	_, span := resource.Tracer.Start(r.Context(), "TReactInfoHandle")
	defer span.End()

	atom := resource.Atoms.ElementByNumber[resource.Number]

	node := &treactorpb.Node{
		Name:      resource.AppName,
		Version:   resource.AppVersion,
		Framework: resource.Framework,
		Request: &treactorpb.TReactorRequest{
			Path:    r.RequestURI,
			Headers: make(map[string]string, len(r.Header)),
		},
		Bonds: nil,
		Atom: &treactorpb.Atom{
			Number: resource.Number,
			Symbol: atom.Symbol,
			Name:   atom.Name,
			Period: &atom.Period,
			Group:  &atom.Group,
		},
	}
	for key, values := range r.Header {
		node.Request.Headers[key] = strings.Join(values, "|")
	}

	bytes, _ := protojson.Marshal(node)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func TReactReactionsHandle(w http.ResponseWriter, r *http.Request) {
	_, span := resource.Tracer.Start(r.Context(), "TReactReactionsHandle")
	defer span.End()

	node := &treactorpb.Node{
		Name:      resource.AppName,
		Version:   resource.AppVersion,
		Framework: resource.Framework,
		Request: &treactorpb.TReactorRequest{
			Path:    r.RequestURI,
			Headers: make(map[string]string, len(r.Header)),
		},
	}
	for key, values := range r.Header {
		node.Request.Headers[key] = strings.Join(values, "|")
	}

	url := r.URL
	molecule := url.Query().Get("molecule")
	if len(molecule) < 3 {
		/*
		   		        axios.get<any, AxiosResponse<Node>>(Config.atomUrl(molecule)).then(
		                  function (result) {
		                      let response: TReactorResponse = {
		                          statusCode: result.status,
		                          statusMessage: result.statusText,
		                          headers: extractHeadersFromResponse(result)
		                      }
		                      let bond: Bond = {
		                          response: response,
		                          node: result.data
		                      }
		                      node.bonds.push(bond)
		                      res.send(node)
		                  }
		              ).catch(
		                  function (result) {
		                      res.status(502)
		                      res.send(result)
		                  }
		              )

		*/
	} else {
		/*
		   axios.get<any, AxiosResponse<Node>>(Config.moleculeUrl(molecule)).then(
		       function (result) {
		           let response: TReactorResponse = {
		               statusCode: result.status,
		               statusMessage: result.statusText,
		               headers: extractHeadersFromResponse(result)
		           }
		           let bond: Bond = {
		               response: response,
		               node: result.data
		           }
		           node.bonds.push(bond)
		           res.send(node)
		       }
		   ).catch(
		       function (result) {
		           res.status(502)
		           res.send(result)
		       }
		   )
		*/
	}
	bytes, _ := protojson.Marshal(node)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func TReactorHealthz(_ http.ResponseWriter, _ *http.Request) {
}

// You can have a catch all tracer on the route, but it's better to instrument the handlers separate
func instrumentedGet(mux *http.ServeMux, route string, handleFunction func(w http.ResponseWriter, r *http.Request)) {
	fullRoute := fmt.Sprintf("%s%s", resource.Base, route)
	mux.Handle(fullRoute, otelhttp.NewHandler(http.HandlerFunc(handleFunction), fmt.Sprintf("GET %s", fullRoute)))
}

func Serve() {
	atoms := element.NewAtoms()

	fmt.Printf("Telemetry Reactor (%s:%s) listening on port %s\n", resource.AppName, resource.AppVersion, resource.Port)
	fmt.Printf("Mode: %s\n", resource.Mode)

	r := http.NewServeMux()
	r.HandleFunc("/healthz", TReactorHealthz)
	instrumentedGet(r, fmt.Sprintf("/nodes/%d/health", resource.Number), TReactorHealthz)
	instrumentedGet(r, fmt.Sprintf("/nodes/%d/info", resource.Number), TReactInfoHandle)
	instrumentedGet(r, "/reactions", TReactSplitHandle)
	for i := 1; i <= resource.MaxBond; i++ {
		instrumentedGet(r, fmt.Sprintf("/bonds/%d", i), TReactBondHandle)
	}
	instrumentedGet(r, "/bonds/n", TReactBondHandle)
	for sym := range atoms.ElementByName {
		instrumentedGet(r, fmt.Sprintf("/atoms/%s", strings.ToLower(sym)), TReactAtomHandle)
	}
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", resource.Port), r))
}
