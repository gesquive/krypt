package cmd

import (
	"github.com/gesquive/cli"
	"github.com/gesquive/krypt/crypto"
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:       "seal [flags] FILE [FILE...]",
	Aliases:   []string{"s"},
	Short:     "Seal unencrypted file(s)",
	Long:      `Seal existing unencrypted files. This command can operate on multiple files at once.`,
	ValidArgs: []string{"FILE"},
	PreRun:    runPreCheck,
	Run:       runEncrypt,
}

func init() {
	RootCmd.AddCommand(encryptCmd)
}

func runEncrypt(cmd *cobra.Command, args []string) {
	if len(args) <= 0 {
		cli.Info("No file to encrypt specified.")
		return
	}
	for _, file := range args {
		cli.Debug("Encrypting %s", file)
		err := encryptFile(password, file)
		if err != nil {
			if _, ok := err.(*crypto.DataIsEncryptedError); ok {
				cli.Error("File is already encrypted, will not encrypt again")
				continue
			}
			cli.Error("Could not encrypt %s", file)
			cli.Debug("%v", err)
		}
	}
}
