package cmd

import (
	"github.com/gesquive/cli"
	"github.com/gesquive/krypt/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// sealCmd represents the encrypt command
var sealCmd = &cobra.Command{
	Use:       "seal [flags] FILE [FILE...]",
	Aliases:   []string{"s"},
	Short:     "Seal unencrypted file(s)",
	Long:      `Seal existing unencrypted files. This command can operate on multiple files at once.`,
	ValidArgs: []string{"FILE"},
	PreRun:    runSealPreRun,
	Run:       runSeal,
}

func init() {
	RootCmd.AddCommand(sealCmd)

	sealCmd.PersistentFlags().StringP("password-file", "p", "",
		"The password file")
	sealCmd.PersistentFlags().StringP("cipher", "i", "AES256",
		"The cipher to encrypt with. Use the list command for a full list.")

	viper.BindEnv("cipher")
	viper.BindEnv("password")
	viper.BindEnv("password-file")
}

func runSealPreRun(cmd *cobra.Command, args []string) {
	viper.BindPFlag("cipher", cmd.PersistentFlags().Lookup("cipher"))
	viper.BindPFlag("password-file", cmd.PersistentFlags().Lookup("password-file"))
}
func runSeal(cmd *cobra.Command, args []string) {
	cipherType := cliGetCipherType()
	password := cliGetPassword()

	if len(args) <= 0 {
		cli.Info("No file to encrypt specified.")
		return
	}
	for _, file := range args {
		cli.Debug("Encrypting %s", file)
		err := encryptFile(cipherType, password, file)
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
