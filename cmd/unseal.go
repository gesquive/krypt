package cmd

import (
	"github.com/gesquive/cli"
	"github.com/gesquive/krypt/crypto"
	"github.com/spf13/cobra"
)

// unsealCmd represents the decrypt command
var unsealCmd = &cobra.Command{
	Use:       "unseal [flags] FILE [FILE...]",
	Aliases:   []string{"u", "unsl"},
	Short:     "Unseal encrypted file(s)",
	Long:      `Unseal existing encrypted files. This command can operate on multiple files at once.`,
	ValidArgs: []string{"FILE"},
	PreRun:    runPreCheck,
	Run:       runUnseal,
}

func init() {
	RootCmd.AddCommand(unsealCmd)
}

func runUnseal(cmd *cobra.Command, args []string) {
	if len(args) <= 0 {
		cli.Info("No file to decrypt specified.")
		return
	}
	for _, file := range args {
		cli.Debug("Decrypting %s", file)
		err := decryptFile(password, file)
		if err != nil {
			if _, ok := err.(*crypto.DataIsNotEncryptedError); ok {
				cli.Error("File is not encrypted, cannot decrypt")
				continue
			}
			cli.Error("Could not decrypt %s", file)
			cli.Debug("%v", err)
		}
	}
}
