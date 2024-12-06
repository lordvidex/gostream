/*
Copyright © 2024 Evans Owamoyo evans.dev99@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/lordvidex/gostream/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var writeConfig bool
var cfg config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gostream",
	Short: "gostream CLI tool",
	Long:  `gostream is a command line tool that creates/configures clients and servers of gostream.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if writeConfig {
			if err := viper.SafeWriteConfigAs("gostream_new.toml"); err != nil {
				return err
			}
			fmt.Println("current config written to ./gostream_new.toml")
		}
		return nil
	},
}

// Execute ...
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/gostream.toml)")
	rootCmd.PersistentFlags().BoolVarP(&writeConfig, "write-config", "w", false, "write out the current config to a toml file")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("toml")
		viper.SetConfigName("gostream")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Fprintln(os.Stderr, "unmarshalling config file:", err)
	}
}
