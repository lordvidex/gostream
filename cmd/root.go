/*
Copyright © 2024 Evans Owamoyo evans.dev99@gmail.com
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lordvidex/gostream/internal/config"
)

var cfgFile string
var writeConfig string
var cfg config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gostream",
	Short: "gostream CLI tool",
	Long:  `gostream is a command line tool that creates/configures clients and servers of gostream.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if len(writeConfig) > 0 {
			fileName := fmt.Sprintf("%s/%s.toml", filepath.Dir(writeConfig), strings.TrimSuffix(filepath.Base(writeConfig), filepath.Ext(writeConfig)))
			if err := viper.SafeWriteConfigAs(fileName); err != nil {
				return err
			}
			fmt.Printf("current config written to %s", fileName)
			os.Exit(0)
		}
		return nil
	},
}

// Execute ...
func Execute() {
	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/gostream.toml)")
	rootCmd.PersistentFlags().StringVarP(&writeConfig, "output-config", "o", "", "current config is written to the provided toml file")
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
