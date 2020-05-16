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
	Args:      VerifyExactFileArgs(1),
	PreRun:    runEditPreRun,
	Run:       runEdit,
}

func init() {
	RootCmd.AddCommand(editCmd)

	editCmd.LocalFlags().StringP("editor", "e", "", "The editor to use")
	editCmd.LocalFlags().StringP("password-file", "p", "",
		"The password file")
	editCmd.LocalFlags().StringP("cipher", "i", "AES256",
		"The cipher to encrypt with. Use the list command for a full list.")

	viper.BindEnv("editor")
	viper.BindEnv("cipher")
	viper.BindEnv("password")
	viper.BindEnv("password-file")
}

func runEditPreRun(cmd *cobra.Command, args []string) {
	viper.BindPFlag("editor", cmd.LocalFlags().Lookup("editor"))
	viper.BindPFlag("cipher", cmd.LocalFlags().Lookup("cipher"))
	viper.BindPFlag("password-file", cmd.LocalFlags().Lookup("password-file"))
}

func runEdit(cmd *cobra.Command, args []string) {
	cipherType := cliGetCipherType()
	password := cliGetPassword()
	editor := cliGetEditor()

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

	newPlainText, err := cliRunFileEdit(editor, origPlainText)
	if err != nil {
		cli.Error("Error while editing file '%s'", file)
		cli.Debug("%v", err)
		return
	}

	if bytes.Compare(origPlainText, newPlainText) == 0 {
		cli.Info("File contents have not changed, not modifying")
		return
	}

	if err := writeCrypt(cipherType, password, file, newPlainText); err != nil {
		cli.Error("Could not encrypt data for file '%s'", file)
		cli.Debug("%v", err)
		return
	}
}
