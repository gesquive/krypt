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
	PreRun:    runPreCheck,
	Run:       runCreate,
}

func init() {
	RootCmd.AddCommand(createCmd)

	createCmd.PersistentFlags().StringP("editor", "e", "", "The editor to use")

	viper.BindEnv("editor")
	viper.BindPFlag("editor", createCmd.PersistentFlags().Lookup("editor"))

}

func runCreate(cmd *cobra.Command, args []string) {
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

	newPlainText, err := getFileEdit(editor, []byte(""))
	if err != nil {
		cli.Error("Error while editing file '%s'", file)
		cli.Debug("%v", err)
		return
	}

	if err := writeCrypt(password, file, newPlainText); err != nil {
		cli.Error("Could not encrypt data for file '%s'", file)
		cli.Debug("%v", err)
		return
	}
}
