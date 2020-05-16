package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/gesquive/cli"
	"github.com/gesquive/krypt/crypto"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

// VerifyMinimumNFileArgs returns an error if there is not at least N args.
func VerifyMinimumNFileArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < n {
			return fmt.Errorf("Not enough files specified, expected at least %d", n)
		}
		return nil
	}
}

// VerifyExactFileArgs returns an error if there are not exactly N args.
func VerifyExactFileArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < n {
			return errors.Errorf("Not enough files specified, expected %d", n)
		} else if len(args) > n {
			return errors.Errorf("Too many files specified, expected %d", n)
		}
		return nil
	}
}

// readFile opens a file and reads the content
func readFile(filePath string) ([]byte, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errors.Wrapf(err, "could not open file to read")
	}

	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read file")
	}
	return contents, nil
}

// writeFile opens a file and write contents in
func writeFile(filePath string, contents []byte) error {
	fileObj, err := os.Create(filePath)
	if err != nil {
		return errors.Wrapf(err, "could not open file to write")
	}
	defer fileObj.Close()

	if _, err = io.Copy(fileObj, bytes.NewReader(contents)); err != nil {
		return errors.Wrapf(err, "could not write to file")
	}
	return nil
}

// readCrypt opens a file, reads it and decrypts the contents
func readCrypt(password string, filePath string) ([]byte, error) {
	var empty []byte
	cipherText, err := readFile(filePath)
	if err != nil {
		return empty, err
	}

	plainText, err := crypto.Decrypt([]byte(password), cipherText)
	if err != nil {
		if derr, ok := err.(*crypto.DataIsNotEncryptedError); ok {
			return empty, derr
		}
		return empty, errors.Wrapf(err, "could not encrypt data")
	}
	return plainText, nil
}

// writeCrypt encrypts the plain text and writes to filePath
func writeCrypt(cipherType crypto.CipherType, password string, filePath string, plainText []byte) error {
	cipherText, err := crypto.Encrypt(cipherType, []byte(password), plainText)
	if err != nil {
		if derr, ok := err.(*crypto.DataIsEncryptedError); ok {
			return derr
		}
		return errors.Wrapf(err, "could not encrypt data")
	}

	if err := writeFile(filePath, cipherText); err != nil {
		return err
	}
	return nil
}

// encryptFile opens a file, encrypts the contents, and writes back the cipher text
func encryptFile(cipherType crypto.CipherType, password string, filePath string) error {
	plainText, err := readFile(filePath)
	if err != nil {
		return err
	}

	err = writeCrypt(cipherType, password, filePath, plainText)
	if err != nil {
		return err
	}
	return nil
}

// decryptFile opens a file, decrypts the contents, and writes back the plain text
func decryptFile(password string, filePath string) error {
	plainText, err := readCrypt(password, filePath)
	if err != nil {
		return err
	}

	if err := writeFile(filePath, plainText); err != nil {
		return err
	}
	return nil
}

func cliGetPassword() string {
	// if a password is provided, use it
	envPassword := strings.TrimSpace(viper.GetString("password"))
	if len(envPassword) > 0 {
		cli.Debug("password src: cli")
		return viper.GetString("password")
	}
	// if a password-file is provided, use the password in it
	passwordFilePath := viper.GetString("password-file")
	if len(passwordFilePath) > 0 {
		if _, err := os.Stat(passwordFilePath); os.IsNotExist(err) {
			cli.Error("password-file: does not exist (\"%s\")", passwordFilePath)
		} else {
			filePassword, err := ioutil.ReadFile(passwordFilePath)
			if err != nil {
				cli.Error("password-file: could not open (\"%s\")", passwordFilePath)
			} else {
				filePassword = bytes.TrimSpace(filePassword)
				if len(filePassword) > 0 {
					cli.Debug("password-file: \"%s\"", passwordFilePath)
					return string(filePassword)
				}
				cli.Error("password-file: file is empty (\"%s\")", passwordFilePath)
			}
		}

	}
	// no password has been provided, kindly pester the user for a valid password
	var userPassword []byte
	for len(userPassword) == 0 {
		fmt.Print("Enter password: ")
		userPassword, _ = terminal.ReadPassword(int(syscall.Stdin))
		fmt.Print("\n")
		userPassword = bytes.TrimSpace(userPassword)
		if len(userPassword) == 0 {
			cli.Error("Password is not long enough")
		}
	}
	return string(userPassword)
}

func cliGetCipherType() crypto.CipherType {
	cipherName := viper.GetString("cipher")
	cipherType, err := crypto.GetCipherTypeByName(cipherName)
	if err != nil || cipherType == crypto.Unknown {
		cli.Fatal("unknown encryption cipher specified")
	}

	cli.Debug("cipher: '%s'", cipherName)
	return cipherType
}

// cliRunFileEdit creates a temporary file and opens it with the given editor.
// 	If the editor successfully returns, the contents of the temporary file are returned.
func cliRunFileEdit(editor string, content []byte) ([]byte, error) {
	// create temp file, w/ permissions
	tmpFile, err := ioutil.TempFile("", "krypt")
	if err != nil {
		return nil, errors.Wrapf(err, "creating tempfile")
	}

	cli.Debug("tmpfile: %s", tmpFile.Name())
	defer func() {
		tmpFile.Close()
		cli.Debug("removing tmpfile")
		os.Remove(tmpFile.Name()) // clean up later
	}()

	// write contents to temp file for editing
	if _, err = tmpFile.Write(content); err != nil {
		return nil, errors.Wrapf(err, "writing to tempfile failed")
	}
	// start editor with file
	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, errors.Wrapf(err, "editor start failed")
	}
	cli.Debug("started editor")

	if err := cmd.Wait(); err != nil {
		// editor failed, just delete the tempfile
		cli.Debug("editor returns failure: %v", err)
	} else {
		// editor success, return contents
		cli.Debug("editor returns success")
		newContent, err := readFile(tmpFile.Name())
		if err == nil {
			return newContent, nil
		}
		return nil, err
	}
	return nil, nil
}

// getEditor gets an editor to use by first checking env variables to see if they are defined,
//	if not, we go through a list of known editors we might be able to use. If any are found, we use them
func cliGetEditor() string {
	cmdEditor := viper.GetString("editor")
	if len(cmdEditor) > 0 {
		cli.Debug("cli editor: '%s'", cmdEditor)
		return cmdEditor
	}
	envEditor := os.Getenv("EDITOR")
	if len(envEditor) > 0 {
		cli.Debug("env editor: '%s'", envEditor)
		return envEditor
	}

	// user didn't supply an editor, try to find some common ones
	knownEditors := []string{"vim", "vi", "nano"}
	for _, knownEditor := range knownEditors {
		path, err := exec.LookPath(knownEditor)
		if err == nil {
			cli.Debug("editor: '%s'", path)
			return path
		}
	}

	cli.Fatal("No editor found, please specify an editor")
	return ""
}
