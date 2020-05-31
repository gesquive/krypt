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
	Short:   "Decrypt and view the contents of a sealed file without editing",
	Long: `This command will decrypt the file to a temporary file and allow you to view the
file without modifying the contents.`,
	ValidArgs: []string{"FILE"},
	Args:      VerifyExactFileArgs(1),
	PreRun:    runViewPreRun,
	Run:       runView,
}

func init() {
	RootCmd.AddCommand(viewCmd)

	viewCmd.PersistentFlags().StringP("editor", "e", "", "The editor to use")
	viewCmd.PersistentFlags().StringP("password-file", "p", "",
		"The password file")

	viper.BindEnv("editor")
	viper.BindEnv("password")
	viper.BindEnv("password-file")
}

func runViewPreRun(cmd *cobra.Command, args []string) {
	viper.BindPFlag("editor", cmd.PersistentFlags().Lookup("editor"))
	viper.BindPFlag("password-file", cmd.PersistentFlags().Lookup("password-file"))
}

func runView(cmd *cobra.Command, args []string) {
	password := cliGetPassword()
	editor := cliGetEditor()

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

	cliRunFileEdit(editor, plainText)
}
