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

var newPassword string

// rekeyCmd represents the rekey command
var rekeyCmd = &cobra.Command{
	Use:       "rekey [flags] FILE",
	Aliases:   []string{"r", "key", "reky"},
	Short:     "Change the password on encrypted file(s)",
	Long:      `Change the password on encrypted file(s). This command can operate on multiple files at once.`,
	ValidArgs: []string{"FILE"},
	PreRun:    runReKeyPreCheck,
	Run:       runReKey,
}

func init() {
	RootCmd.AddCommand(rekeyCmd)

	rekeyCmd.PersistentFlags().StringP("new-password-file", "n", "",
		"The new password file")

	viper.BindEnv("new-password")
	viper.BindEnv("new-password-file")

	viper.BindPFlag("new-password-file", rekeyCmd.PersistentFlags().Lookup("new-password-file"))
}

func runReKeyPreCheck(cmd *cobra.Command, args []string) {
	runPreCheck(cmd, args)
	newPassword = getNewPassword()
}

func runReKey(cmd *cobra.Command, args []string) {
	cli.Info("pass: %s %s", password, newPassword)
	if len(args) <= 0 {
		cli.Info("No file to re-key specified.")
		return
	}
	for _, file := range args {
		cli.Debug("rekey %s", file)
		plainText, err := readCrypt(password, file)
		if err != nil {
			if _, ok := err.(*crypto.DataIsNotEncryptedError); ok {
				cli.Error("File is not encrypted, cannot decrypt")
				continue
			}
			cli.Error("Could not decrypt %s", file)
			cli.Debug("%v", err)
		}

		err = writeCrypt(newPassword, file, plainText)
		if err != nil {
			cli.Error("Could not write to %s", file)
			cli.Debug("%v", err)
		}
	}
}

func getNewPassword() string {
	// if a password is provided, use it
	envPassword := strings.TrimSpace(viper.GetString("new-password"))
	if len(envPassword) > 0 {
		cli.Debug("Found new password in environment variables")
		return viper.GetString("new-password")
	}
	// if a password-file is provided, use the password in it
	passwordFilePath := viper.GetString("new-password-file")
	if len(passwordFilePath) > 0 {
		if _, err := os.Stat(passwordFilePath); !os.IsNotExist(err) {
			cli.Error("new-password-file: \"%s\" does not exist")
		} else {
			filePassword, err := ioutil.ReadFile(passwordFilePath)
			if err != nil {
				cli.Error("new-password-file: could not open")
			} else {
				filePassword = bytes.TrimSpace(filePassword)
				if len(filePassword) > 0 {
					cli.Debug("Got new password from new-password-file")
					return string(filePassword)
				}
				cli.Error("new-password-file: file is empty")
			}
		}

	}
	// no password has been provided, kindly pester the user for a valid password
	var userPassword []byte
	for len(userPassword) == 0 {
		fmt.Print("Enter new password: ")
		userPassword, _ = terminal.ReadPassword(int(syscall.Stdin))
		fmt.Print("\n")
		userPassword = bytes.TrimSpace(userPassword)
		if len(userPassword) == 0 {
			cli.Error("new password is not long enough")
		}
	}
	return string(userPassword)
}
