package resource

import (
	"github.com/treactor/treactor-go/pkg/element"
)

var (
	Logger RLogger
)

var Atoms *element.Atoms

func Init() {
	initTelemetry()
	clientInit()
	Logger = NewSLogger("")
	Atoms = element.NewAtoms()
}
