# Telemetry Reactor for GoLang

TReactor is a microservice designed to test and experiment with observability of microservices. You can play with it
on your own machine in `local` more, but it gets interesting when you deploy it on a Kubernetes cluster or an Istio mesh.

In cluster mode you can have a bit over 120 microservices (atoms in the mendeleev table). You control `treactor` by
giving it some interesting `molecules`.

### Motivation

I created this microservice to let me inspect what happens inside the kubernetes cluster or a mesh. I've noticed that
not all the scenario's and products worked that well together, certainly when you go a bit beyond the happy path. This
microservice will enable me to *make reproducibles*. I hope it's also helpful for other people to *learn*  technologies
like tracing, kubernetes networking, logging, istio, etc...

I will start using (and extending) treactor for bug reports, articles and talks. Let's see if this will get
as popular as httpbin is ;).

### Features

* Create cascading microservice calls from one origin call
* Configure tracing and logging
* Inspect headers, every level in the call hierarchy

### Example

To see what treactor does lets take a very simple example. TReactor works by splitting molecules:

`[[H]]^2[O]`

Everything between the brackets `[]` will be a call to the next microservice. The brackets can be prefixed with a number
and an optional parameter (`s` or `p`), this tells treactor how many times the service needs to be called and how (
sequential or parallel). Multiple calls can be make by appending them using `^` (sequential) or `*` (parallel).

Depending the content of the bracket the call will be different. If treactor detects an atom a call to the corresponding
atom service will be made. But if treactor detects another sub-molecule it calls the next bond and apply the same
logic till only atoms are left. So the example above will result in:

`http://treactor-api/treact/split?molecule=[[H]]^2[O]`

calling

* `http://bond-1/treact/split?molecule=[h]`
* `http://atom-o/treact/atom/o?atom=o`
* `http://atom-o/treact/atom/o?atom=o`

*bond* will split the molecule [H] (ok, this looks strange, but each bracket is a layer) into it's atoms, in this
case only 1 `H`:

* `http://atom-h/treact/atom?atom=H`

Try the local installation, to see how it looks in the trace (this will make it more clear).

## Installation

### Pre-Requirement

Create a directory `tmp` and `work` in this repo. Don't worry, they are in the `.gitignore` so you do not accidentally
check them in.

### Local

Build the `treactor` from source

`go install ./cmd/treactor/`

Set the environmental variables. Replace `project-name` project with your own

```
export PORT=3330
export TREACTOR_NAME=reactor-api
export TREACTOR_VERSION=1
export TREACTOR_DEBUG=1
export TREACTOR_PROFILE=0
export TREACTOR_MODE=local
```

Fire up the treactor

`treactor`

And test it by calling:

http://localhost:3330/treact/split?molecule=[[H]]^2[O]

Go to the Cloud Console, select *Trace*.

### Kubernetes

*Not yet fully tested/supported*


`go install ./cmd/trprep/`

`kubectl label namespace default istio-injection=disabled --overwrite`

### Istio

*Not yet fully tested/supported*

`kubectl label namespace default istio-injection=enabled --overwrite`

## Specification

### Reactor Prepare (trprep)

Prepares the Kubernetes files, from the templates

### Environment Variables

NAME | Description | Default
---- | ----------- | -------
PORT | Port |
SERVICE_NAME | Application name | treactor
SERVICE_VERSION | Application version | 0.0.0
TREACTOR_MODE | Reactor mode (local, k8s) | local
TREACTOR_TRACE_PROPAGATION | OpenTelemetry propagator (w3c)  | w3c

### Molecule spec

```
S    [H,x=1,y=2,z=3]*2[O,x=1,y=2,z=3],x=1,y=2,z=3
A     H,x=1,y=2,z=3
A                      O,x=1,y=2,z=3
```


```
S    2[5[Ur,log:1,xyz:4]^5[C,log:1,xyz:4]],x:1,y:2
O1     5[Ur,log:1,xyz:4]^5[C,log:1,xyz:4]
A        Ur,log:1,xyz:4
A                          C,log:1,xyz:4
```

## Test Scenarios

https://github.com/wg/wrk

`wrk -t12 -c400 -d30s "http://<yourip>/treact/split?molecule=6[C]^12[H]^6[O]"`

glucose
6[C]^12[H]^6[O]

https://www.sigmaaldrich.com/catalog/product/aldrich/375756?lang=en&region=BE
Baicalin hydrate

C21H18O11 · xH2O

env GOOS=linux GOARCH=amd64 go build cmd/treactor/main.go