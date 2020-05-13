package main

//go:generate go run generate_manpages.go

import (
	"github.com/gesquive/krypt/cmd"
)

func main() {
	cmd.Execute()
}
