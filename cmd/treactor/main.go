package main

import (
	"github.com/treactor/treactor-go/pkg/resource"
	"github.com/treactor/treactor-go/pkg/treact"
)

func main() {
	resource.Configure()
	resource.Init()
	treact.Serve()
}
