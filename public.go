package generatorgit

import (
	"github.com/StephanHCB/go-generator-git/api"
	"github.com/StephanHCB/go-generator-git/internal/implementation"
)

var Instance api.GitApi

func init() {
	Instance = &implementation.GitGeneratorImpl{}
}
