package cmd

import (
	"github.com/gesquive/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:     "create [flags] FILE",
	Aliases: []string{"c"},
	Short:   "Create a new encrypted text file",
	Long: `This command will create a new file and allow you to edit the file using the
defined editor. After editing the file, the contents will be encrypted to the defined file.`,
	ValidArgs: []string{"FILE"},
	Args:      VerifyExactFileArgs(1),
	PreRun:    runCreatePreRun,
	Run:       runCreate,
}

func init() {
	RootCmd.AddCommand(createCmd)

	createCmd.LocalFlags().StringP("editor", "e", "", "The editor to use")
	createCmd.LocalFlags().StringP("password-file", "p", "",
		"The password file")
	createCmd.LocalFlags().StringP("cipher", "i", "AES256",
		"The cipher to encrypt with. Use the list command for a full list.")

	viper.BindEnv("editor")
	viper.BindEnv("cipher")
	viper.BindEnv("password")
	viper.BindEnv("password-file")
}

func runCreatePreRun(cmd *cobra.Command, args []string) {
	viper.BindPFlag("editor", cmd.LocalFlags().Lookup("editor"))
	viper.BindPFlag("cipher", cmd.LocalFlags().Lookup("cipher"))
	viper.BindPFlag("password-file", cmd.LocalFlags().Lookup("password-file"))
}

func runCreate(cmd *cobra.Command, args []string) {
	cipherType := cliGetCipherType()
	password := cliGetPassword()
	editor := cliGetEditor()

	file := args[0]

	newPlainText, err := cliRunFileEdit(editor, []byte(""))
	if err != nil {
		cli.Error("Error while editing file '%s'", file)
		cli.Debug("%v", err)
		return
	}

	if err := writeCrypt(cipherType, password, file, newPlainText); err != nil {
		cli.Error("Could not encrypt data for file '%s'", file)
		cli.Debug("%v", err)
		return
	}
}
