package resource

import (
	"github.com/treactor/treactor-go/pkg/chem"
)

var (
	Logger RLogger
)

var Atoms *chem.Atoms

func Init() {
	initTelemetry()
	clientInit()
	Logger = NewSLogger("")
	Atoms = chem.NewAtoms()
}
