// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gesquive/krypt/cmd"
	"github.com/spf13/cobra/doc"
)

var destinationPath = "manpages"

func main() {
	if _, err := os.Stat(destinationPath); os.IsNotExist(err) {
		os.MkdirAll(destinationPath, 0755)
	}

	header := &doc.GenManHeader{
		Title:   "KRYPT",
		Section: "1",
		Source:  "krypt",
		Manual:  "krypt utils",
	}
	cmd.RootCmd.DisableAutoGenTag = true
	fmt.Printf("generating manpages for krypt\n")

	if err := doc.GenManTree(cmd.RootCmd, header, destinationPath); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(2)
	}

	//Remove all of the double blank lines from output docs
	err := filepath.Walk(destinationPath, func(path string, f os.FileInfo, err error) error {
		stripFile(path)
		return nil
	})

	if err != nil {
		fmt.Printf("Could not clean up all the files\n")
		fmt.Printf("%s", err)
	}
}

func stripFile(path string) error {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	regex, err := regexp.Compile("\n{2,}")
	if err != nil {
		return err
	}
	output := regex.ReplaceAllString(string(input), "\n")

	err = ioutil.WriteFile(path, []byte(output), 0644)
	if err != nil {
		return err
	}
	return nil
}
