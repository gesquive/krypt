package cmd

import (
	"github.com/gesquive/cli"
	"github.com/gesquive/krypt/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// unsealCmd represents the decrypt command
var unsealCmd = &cobra.Command{
	Use:       "unseal [flags] FILE [FILE...]",
	Aliases:   []string{"u", "unsl"},
	Short:     "Unseal encrypted file(s)",
	Long:      `Unseal existing encrypted files. This command can operate on multiple files at once.`,
	ValidArgs: []string{"FILE"},
	Args:      VerifyMinimumNFileArgs(1),
	PreRun:    runUnsealPreRun,
	Run:       runUnseal,
}

func init() {
	RootCmd.AddCommand(unsealCmd)

	unsealCmd.PersistentFlags().StringP("password-file", "p", "",
		"The password file")

	viper.BindEnv("password")
	viper.BindEnv("password-file")
}

func runUnsealPreRun(cmd *cobra.Command, args []string) {
	viper.BindPFlag("password-file", cmd.PersistentFlags().Lookup("password-file"))
}

func runUnseal(cmd *cobra.Command, args []string) {
	password := cliGetPassword()

	for _, file := range args {
		// TODO: use glob to expand file paths
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
