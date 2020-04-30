package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/gesquive/cli"
	"github.com/gesquive/krypt/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	buildVersion = "v0.1.0-dev"
	buildCommit  = ""
	buildDate    = ""
)

var cfgFile string
var displayVersion string

var debug bool
var showVersion bool

var password string
var cipherType crypto.CipherType

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:              "krypt [flags] command",
	Short:            "Encrypt or Decrypt files",
	Long:             `Encrypt or Decrypt files with various ciphers`,
	PersistentPreRun: preRun,
	Hidden:           true,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	RootCmd.SetHelpTemplate(fmt.Sprintf("%s\nVersion:\n  github.com/gesquive/krypt %s\n",
		RootCmd.HelpTemplate(), buildVersion))
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "",
		"Path to a specific config file (default \"./config.yml\")")
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false,
		"Write debug messages to console")
	RootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "V", false,
		"Show the version info and exit")
	RootCmd.PersistentFlags().StringP("password-file", "p", "",
		"The password file")
	RootCmd.PersistentFlags().StringP("cipher", "y", "AES256",
		"The cipher to en/decrypt with. Use the list command for a full list.")

	RootCmd.PersistentFlags().MarkHidden("debug")

	viper.SetEnvPrefix("krypt")
	viper.AutomaticEnv()

	viper.BindEnv("cipher")
	viper.BindEnv("password")
	viper.BindEnv("password-file")

	viper.BindPFlag("cipher", RootCmd.PersistentFlags().Lookup("cipher"))
	viper.BindPFlag("password-file", RootCmd.PersistentFlags().Lookup("password-file"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfgFile := viper.GetString("config")
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")              // name of config file (without extension)
		viper.AddConfigPath(".")                   // add current directory as first search path
		viper.AddConfigPath("$HOME/.config/krypt") // add home directory to search path
		viper.AddConfigPath("/etc/krypt")          // add etc to search path
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if !showVersion {
			if !strings.Contains(err.Error(), "Not Found") {
				fmt.Printf("Error opening config: %s\n", err)
			}
		}
	}
}
func preRun(cmd *cobra.Command, args []string) {
	if showVersion {
		fmt.Printf("github.com/gesquive/krypt\n")
		fmt.Printf(" Version:    %s\n", buildVersion)
		if len(buildCommit) > 6 {
			fmt.Printf(" Git Commit: %s\n", buildCommit[:7])
		}
		if buildDate != "" {
			fmt.Printf(" Build Date: %s\n", buildDate)
		}
		fmt.Printf(" Go Version: %s\n", runtime.Version())
		fmt.Printf(" OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}
	if debug {
		cli.SetPrintLevel(cli.LevelDebug)
	}
	cli.Debug("Running with debug turned on")
}

func runPreCheck(cmd *cobra.Command, args []string) {
	cli.Debug("config:", viper.ConfigFileUsed())

	cipherName := viper.GetString("cipher")
	cipherType = crypto.GetCipherTypeByName(cipherName)
	if cipherType == crypto.Unknown {
		cli.Fatal("Unknown encryption cipher specified")
	}
	cli.Debug("Using cipher '%s'", cipherName)

	password = getPassword()
}

func getPassword() string {
	// if a password is provided, use it
	envPassword := strings.TrimSpace(viper.GetString("password"))
	if len(envPassword) > 0 {
		cli.Debug("Found password in environment variables")
		return viper.GetString("password")
	}
	// if a password-file is provided, use the password in it
	passwordFilePath := viper.GetString("password-file")
	if len(passwordFilePath) > 0 {
		if _, err := os.Stat(passwordFilePath); !os.IsNotExist(err) {
			cli.Error("password-file: \"%s\" does not exist")
		} else {
			filePassword, err := ioutil.ReadFile(passwordFilePath)
			if err != nil {
				cli.Error("password-file: could not open")
			} else {
				filePassword = bytes.TrimSpace(filePassword)
				if len(filePassword) > 0 {
					cli.Debug("got password from password-file")
					return string(filePassword)
				}
				cli.Error("password-file: file is empty")
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
