package execute

import (
	"fmt"
	treactorpb "github.com/treactor/treactor-go/io/treactor/v1alpha"
	"github.com/treactor/treactor-go/pkg/resource"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/otel"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/encoding/protojson"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

type Plan interface {
	String() string
	Execute(ctx context.Context, channel chan *treactorpb.Bond)
	Calls() int
}

type Block struct {
	times int
	mode  string
	Block string
	KV    map[string]string
}

func (o *Block) isAtom() bool {
	return unicode.IsLetter(rune(o.Block[0]))
}

func (o *Block) callElement(ctx context.Context, wg *sync.WaitGroup, channel chan  *treactorpb.Bond) {
	defer wg.Done()
	// If REACTOR_TRACE_INTERNAL=1 add internal spans
	ctx, span := otel.Tracer("").Start(ctx, "Block [callElement]")
	defer span.End()
	CallElementResource(ctx, channel, o.Block)
}

func (o *Block) callBond(ctx context.Context, wg *sync.WaitGroup, channel chan  *treactorpb.Bond) {
	defer wg.Done()
	// If REACTOR_TRACE_INTERNAL=1 add internal spans
	ctx, span := otel.Tracer("").Start(ctx, "Block [callBond]")
	defer span.End()
	CallBondResource(ctx, channel, o.Block)
}

func (o *Block) Execute(ctx context.Context, channel chan  *treactorpb.Bond) {
	// If REACTOR_TRACE_INTERNAL=1 add internal spans
	ctx, span := otel.Tracer("").Start(ctx, "Execute Block")
	defer span.End()
	wg := sync.WaitGroup{}
	wg.Add(o.times)
	if o.mode == "s" {
		for i := 1; i <= o.times; i++ {
			if o.isAtom() {
				o.callElement(ctx, &wg, channel)
			} else {
				o.callBond(ctx, &wg, channel)
			}
		}
	} else if o.mode == "p" {
		for i := 1; i <= o.times; i++ {
			if o.isAtom() {
				go o.callElement(ctx, &wg, channel)
			} else {
				go o.callBond(ctx, &wg, channel)
			}
		}
	} else {
		// TODO ERR
	}
	wg.Wait()
}

func (o *Block) Calls() int {
	return o.times
}

func (o *Block) String() string {
	s := strconv.Itoa(o.times) + o.mode + "[" + o.Block + "]"
	for k, v := range o.KV {
		s += "," + k + ":" + v
	}
	return s
}

type Operator struct {
	left    Plan
	right   Plan
	operand Token
}

func (o *Operator) Execute(ctx context.Context, channel chan  *treactorpb.Bond) {
	// If REACTOR_TRACE_INTERNAL=1 add internal spans
	ctx, span := otel.Tracer("").Start(ctx, "Execute Operator")
	defer span.End()
	wg := sync.WaitGroup{}
	wg.Add(2)
	if o.operand == PLUS {
		o.execute(ctx, &wg, channel, o.left)
		o.execute(ctx, &wg, channel, o.right)
	} else if o.operand == MULTIPLY {
		go o.execute(ctx, &wg, channel, o.left)
		go o.execute(ctx, &wg, channel, o.right)
	} else {
		// TODO ERR
	}
	wg.Wait()
}

func (o *Operator) Calls() int {
	return o.left.Calls() + o.right.Calls()
}

func (o *Operator) execute(ctx context.Context, wg *sync.WaitGroup, channel chan  *treactorpb.Bond, plan Plan) {
	defer wg.Done()
	ctx, span := otel.Tracer("").Start(ctx, "Operator [execute]")
	defer span.End()
	plan.Execute(ctx, channel)
}

func (o *Operator) String() string {
	if o.operand == MULTIPLY {
		return o.left.String() + "*" + o.right.String()
	}
	return o.left.String() + "^" + o.right.String()
}

func CallBondResource(context context.Context, channel chan  *treactorpb.Bond, molecule string) {
	next := resource.NextBond()
	var url string
	if resource.IsLocalMode() {
		url = fmt.Sprintf("http://localhost:%s%s/bond/%s?molecule=%s", resource.Port, resource.Base, next, molecule)
	} else {
		url = fmt.Sprintf("http://bond-%s%s/bond/%s?molecule=%s", next, resource.Base, next, molecule)
	}
	context = httptrace.WithClientTrace(context, otelhttptrace.NewClientTrace(context))
	req, _ := http.NewRequestWithContext(context, "GET", url, nil)
	ra, err := resource.HttpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ra.Body.Close()

	bodyBytes, err := ioutil.ReadAll(ra.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var node treactorpb.Node
	protojson.Unmarshal(bodyBytes, &node)

	var bond = &treactorpb.Bond{
		Response: nil,
		Node:     &node,
	}
	channel <- bond
}

func CallElementResource(context context.Context, channel chan  *treactorpb.Bond, symbol string) {
	full := symbol
	symbol = strings.Split(full, ",")[0]
	var url string
	if resource.IsLocalMode() {
		url = fmt.Sprintf("http://localhost:%s%s/atom/%s?symbol=%s", resource.Port, resource.Base, symbol, full)
	} else {
		url = fmt.Sprintf("http://atom-%s%s/atom/%s?symbol=%s", strings.ToLower(symbol), resource.Base, strings.ToLower(symbol), full)
	}
	context = httptrace.WithClientTrace(context, otelhttptrace.NewClientTrace(context))
	req, _ := http.NewRequestWithContext(context, "GET", url, nil)
	ra, err := resource.HttpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ra.Body.Close()

	bodyBytes, err := ioutil.ReadAll(ra.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var node treactorpb.Node
	protojson.Unmarshal(bodyBytes, &node)

	var bond = &treactorpb.Bond{
		Response: nil,
		Node:     &node,
	}
	channel <- bond

}
