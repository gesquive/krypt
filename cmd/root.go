package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/gesquive/cli"
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
	RootCmd.PersistentFlags().MarkHidden("debug")

	viper.SetEnvPrefix("krypt")
	viper.AutomaticEnv()
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
				cli.Error("Error opening config: %s", err)
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
	cli.Debug("running with debug turned on")
	cli.Debug("config: %s", viper.ConfigFileUsed())
}
