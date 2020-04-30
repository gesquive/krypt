package cmd

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/gesquive/cli"
	"github.com/gesquive/krypt/crypto"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

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
func writeCrypt(password string, filePath string, plainText []byte) error {
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
func encryptFile(password string, filePath string) error {
	plainText, err := readFile(filePath)
	if err != nil {
		return err
	}

	err = writeCrypt(password, filePath, plainText)
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

// getFileEdit creates a temporary file and opens it with the given editor.
// 	If the editor successfully returns, the contents of the temporary file are returned.
func getFileEdit(editor string, content []byte) ([]byte, error) {
	// create temp file, w/ permissions
	tmpFile, err := ioutil.TempFile("", "krypt")
	if err != nil {
		return nil, errors.Wrapf(err, "creating tempfile")
	}

	cli.Debug("using tmpfile: %s", tmpFile.Name())
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
	cli.Debug("waiting for command to finish.")

	if err := cmd.Wait(); err != nil {
		// editor failed, just delete the tempfile
		cli.Debug("editor failed with: %v", err)
	} else {
		// editor success, return contents
		cli.Debug("editor reports success")
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
func getEditor() string {
	cmdEditor := viper.GetString("editor")
	if len(cmdEditor) > 0 {
		return cmdEditor
	}
	envEditor := os.Getenv("EDITOR")
	if len(envEditor) > 0 {
		return envEditor
	}

	// user didn't supply an editor, try to find some common ones
	knownEditors := []string{"vim", "vi", "nano"}
	for _, knownEditor := range knownEditors {
		path, err := exec.LookPath(knownEditor)
		if err == nil {
			return path
		}
	}

	return ""
}
