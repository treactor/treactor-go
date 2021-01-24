# Implementation Nodes

This is a scratch-pad for ideas, links and whatever that should not belong in the readme.


## Test Scenarios

https://github.com/wg/wrk

`wrk -t12 -c400 -d30s "http://<yourip>/treact/split?molecule=6[C]^12[H]^6[O]"`

glucose
6[C]^12[H]^6[O]

https://www.sigmaaldrich.com/catalog/product/aldrich/375756?lang=en&region=BE
Baicalin hydrate

C21H18O11 · xH2O

env GOOS=linux GOARCH=amd64 go build cmd/treactor/main.go



```shell
protoc --go_out=. \
    --proto_path ../treactor-proto \
    --experimental_allow_proto3_optional \
    io/treactor/v1alpha/atom.proto \
    io/treactor/v1alpha/node.proto
```
