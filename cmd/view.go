package cmd

import (
	"github.com/gesquive/cli"
	"github.com/gesquive/krypt/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:     "view [flags] FILE",
	Aliases: []string{"v"},
	Short:   "Decrypt and view the contents of an encrypted file without editing",
	Long: `This command will decrypt the file to a temporary file and allow you to view the
file without modifying the contents.`,
	ValidArgs: []string{"FILE"},
	PreRun:    runPreCheck,
	Run:       runView,
}

func init() {
	RootCmd.AddCommand(viewCmd)

	viewCmd.PersistentFlags().StringP("editor", "e", "", "The editor to use")

	viper.BindEnv("editor")
	viper.BindPFlag("editor", viewCmd.PersistentFlags().Lookup("editor"))

}

func runView(cmd *cobra.Command, args []string) {
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
	plainText, err := readCrypt(password, file)
	if err != nil {
		if _, ok := err.(*crypto.DataIsNotEncryptedError); ok {
			cli.Error("File is not encrypted, cannot decrypt")
			return
		}
		cli.Error("Could not decrypt file '%s'", file)
		cli.Debug("%v", err)
		return
	}

	getFileEdit(editor, plainText)
}
