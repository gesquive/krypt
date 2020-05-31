package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"github.com/gesquive/cli"
	"github.com/gesquive/krypt/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

// resealCmd represents the reseal command
var resealCmd = &cobra.Command{
	Use:       "reseal [flags] FILE",
	Aliases:   []string{"r", "resl"},
	Short:     "Change the password/cipher on encrypted file(s)",
	Long:      `Change the password/cipher on encrypted file(s). This command can operate on multiple files at once.`,
	ValidArgs: []string{"FILE"},
	Args:      VerifyMinimumNFileArgs(1),
	PreRun:    runResealPreRun,
	Run:       runReseal,
}

func init() {
	RootCmd.AddCommand(resealCmd)

	resealCmd.PersistentFlags().StringP("cipher", "i", "AES256",
		"The cipher to encrypt with. Use the list command for a full list.")
	resealCmd.PersistentFlags().StringP("password-file", "p", "",
		"The password file to encrypt with.")
	resealCmd.PersistentFlags().StringP("old-password-file", "o", "",
		"The old password file to decrypt with.")

	viper.BindEnv("cipher")
	viper.BindEnv("password")
	viper.BindEnv("password-file")
	viper.BindEnv("old-password")
	viper.BindEnv("old-password-file")

}

func runResealPreRun(cmd *cobra.Command, args []string) {
	viper.BindPFlag("cipher", cmd.PersistentFlags().Lookup("cipher"))
	viper.BindPFlag("password-file", cmd.PersistentFlags().Lookup("password-file"))
	viper.BindPFlag("old-password-file", cmd.PersistentFlags().Lookup("old-password-file"))
}

func runReseal(cmd *cobra.Command, args []string) {
	cipherType := cliGetCipherType()
	oldPassword := cliGetOldPassword()
	password := cliGetPassword()

	for _, file := range args {
		cli.Debug("reseal %s", file)
		plainText, err := readCrypt(oldPassword, file)
		if err != nil {
			if _, ok := err.(*crypto.DataIsNotEncryptedError); ok {
				cli.Error("File is not encrypted, cannot decrypt")
				continue
			}
			cli.Error("Could not decrypt %s", file)
			cli.Debug("%v", err)
		}

		err = writeCrypt(cipherType, password, file, plainText)
		if err != nil {
			cli.Error("Could not write to %s", file)
			cli.Debug("%v", err)
		}
	}
}

func cliGetOldPassword() string {
	// if a password is provided, use it
	envPassword := strings.TrimSpace(viper.GetString("old-password"))
	if len(envPassword) > 0 {
		cli.Debug("old-password src: cli")
		return viper.GetString("old-password")
	}
	// if a password-file is provided, use the password in it
	passwordFilePath := viper.GetString("old-password-file")
	if len(passwordFilePath) > 0 {
		if _, err := os.Stat(passwordFilePath); !os.IsNotExist(err) {
			cli.Error("old-password-file: does not exist (\"%s\")", passwordFilePath)
		} else {
			filePassword, err := ioutil.ReadFile(passwordFilePath)
			if err != nil {
				cli.Error("old-password-file: could not open (\"%s\")", passwordFilePath)
			} else {
				filePassword = bytes.TrimSpace(filePassword)
				if len(filePassword) > 0 {
					cli.Debug("old-password-file: (\"%s\")", passwordFilePath)
					return string(filePassword)
				}
				cli.Error("old-password-file: file is empty (\"%s\")", passwordFilePath)
			}
		}

	}
	// no password has been provided, kindly pester the user for a valid password
	var userPassword []byte
	for len(userPassword) == 0 {
		fmt.Print("Enter old password: ")
		userPassword, _ = terminal.ReadPassword(int(syscall.Stdin))
		fmt.Print("\n")
		userPassword = bytes.TrimSpace(userPassword)
		if len(userPassword) == 0 {
			cli.Error("old-password is not long enough")
		}
	}
	return string(userPassword)
}
