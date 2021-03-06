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
	Args:      VerifyMinimumNFileArgs(1),
	PreRun:    runSealPreRun,
	Run:       runSeal,
}

func init() {
	RootCmd.AddCommand(sealCmd)

	sealCmd.PersistentFlags().StringP("password-file", "p", "",
		"The password file")
	sealCmd.PersistentFlags().StringP("cipher", "i", "AES256",
		"The cipher to encrypt with. Use the list command for a full list.")
	sealCmd.PersistentFlags().BoolP("encode-text", "t", false,
		"encode the output in base64")

	viper.BindEnv("cipher")
	viper.BindEnv("password")
	viper.BindEnv("password-file")
	viper.BindEnv("encode-text")
}

func runSealPreRun(cmd *cobra.Command, args []string) {
	viper.BindPFlag("cipher", cmd.PersistentFlags().Lookup("cipher"))
	viper.BindPFlag("password-file", cmd.PersistentFlags().Lookup("password-file"))
	viper.BindPFlag("encode-text", cmd.PersistentFlags().Lookup("encode-text"))
}
func runSeal(cmd *cobra.Command, args []string) {
	cipherType := cliGetCipherType()
	password := cliGetPassword()
	encodeText := viper.GetBool("encode-text")

	for _, file := range args {
		cli.Debug("Encrypting %s", file)
		err := encryptFile(cipherType, password, file, encodeText)
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
