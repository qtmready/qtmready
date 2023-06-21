package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"go.breu.io/ctrlplane/tools/linters"
)

func main() {
	singlechecker.Main(linters.KVLogger)
}
