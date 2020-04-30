package cmd

import (
	"bytes"

	"github.com/gesquive/cli"
	"github.com/gesquive/krypt/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:     "edit [flags] FILE",
	Aliases: []string{"e"},
	Short:   "Decrypt, edit and encrypt an encrypted file",
	Long: `This command will decrypt the file to a temporary file and allow you to edit the
file using the defined editor. Once editing is done, it will encrypt the contents back to the
original file.`,
	ValidArgs: []string{"FILE"},
	PreRun:    runPreCheck,
	Run:       runEdit,
}

func init() {
	RootCmd.AddCommand(editCmd)

	editCmd.PersistentFlags().StringP("editor", "e", "", "The editor to use")

	viper.BindEnv("editor")
	viper.BindPFlag("editor", editCmd.PersistentFlags().Lookup("editor"))
}

func runEdit(cmd *cobra.Command, args []string) {
	editor := getEditor()
	if len(editor) == 0 {
		cli.Error("No editor found, please specify an editor")
		return
	}
	cli.Debug("Using '%s' as editor", editor)

	if len(args) <= 0 {
		cmd.Usage()
		return
	}
	file := args[0]
	origPlainText, err := readCrypt(password, file)
	if err != nil {
		if _, ok := err.(*crypto.DataIsNotEncryptedError); ok {
			cli.Error("File is not encrypted, cannot decrypt")
			return
		}
		cli.Error("Could not decrypt file '%s'", file)
		cli.Debug("%v", err)
		return
	}

	newPlainText, err := getFileEdit(editor, origPlainText)
	if err != nil {
		cli.Error("Error while editing file '%s'", file)
		cli.Debug("%v", err)
		return
	}

	if bytes.Compare(origPlainText, newPlainText) == 0 {
		cli.Info("File contents have not changed, not modifying")
		return
	}

	if err := writeCrypt(password, file, newPlainText); err != nil {
		cli.Error("Could not encrypt data for file '%s'", file)
		cli.Debug("%v", err)
		return
	}
}
