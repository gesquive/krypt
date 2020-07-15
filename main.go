package main

//go:generate go run generate_manpages.go

import (
	"github.com/gesquive/krypt/cmd"
)

// current build info
var (
	BuildVersion = "v1.1.0-dev"
	BuildCommit  = ""
	BuildDate    = ""
)

func main() {
	cmd.BuildVersion = BuildVersion
	cmd.BuildCommit = BuildCommit
	cmd.BuildDate = BuildDate
	cmd.Execute()
}
