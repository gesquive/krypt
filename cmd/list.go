package cmd

import (
	"github.com/gesquive/cli"
	"github.com/gesquive/krypt/crypto"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List the available cipher methods",
	Long:    `List the name and description of all the available cipher methods`,
	Run:     runList,
}

func init() {
	RootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {

	cli.Info("Supported Ciphers:")
	cipherList := crypto.GetCipherList()
	for _, cipher := range cipherList {
		cli.Info("%10s  %40s", cipher.GetName(), cipher.GetDescription())
	}
}
