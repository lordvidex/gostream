/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lordvidex/gostream/internal/cmd/client/cli"
	"github.com/lordvidex/gostream/internal/cmd/client/tui"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "gostream client commands",
}

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.PersistentFlags().StringSlice("servers", nil, "list of servers that can be connected to, in the cluster")
	clientCmd.PersistentFlags().StringSliceP("entities", "e", []string{"all"}, "entities to watch separated by comma. Valid values include: `all`, `pet`, `user`")
	clientCmd.PersistentFlags().DurationP("timeout", "t", time.Minute*5, "timeout for connecting to the server")

	_ = viper.BindPFlag("client.servers", clientCmd.Flags().Lookup("servers"))
	_ = viper.BindPFlag("client.entities", clientCmd.Flags().Lookup("entities"))
	_ = viper.BindPFlag("client.connection_timeout", clientCmd.Flags().Lookup("timeout"))
}

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "connects to any streaming server and listens for events",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cli.New(cfg.Client).Watch(cmd.Context()); err != nil {
			fmt.Println("err: ", err)
		}
		fmt.Println("client finished")
	},
}

func init() {
	clientCmd.AddCommand(watchCmd)

	watchCmd.Flags().Int("connections", 1, "number of parallel connections to make with the server. This is strictly for testing purpose, ideally only one connection should be made to server.")
	watchCmd.Flags().StringP("name", "n", uuid.New().String(), "client name")
	watchCmd.Flags().Bool("dry-run", false, "print configs and exit")

	_ = viper.BindPFlag("client.connections", watchCmd.Flags().Lookup("connections"))
	_ = viper.BindPFlag("client.name", watchCmd.Flags().Lookup("name"))
	_ = viper.BindPFlag("client.dry_run", watchCmd.Flags().Lookup("dry-run"))
}

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "create single client terminal user interface",
	Long: `displays a complete Terminal User Interface for performing all functionalities of the client which includes:
	- watching streams from the server
	- creating and editing entities`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := tui.New(cmd.Context(), cfg.Client).Run(); err != nil {
			fmt.Println("err: ", err)
		}
	},
}

func init() {
	clientCmd.AddCommand(tuiCmd)

	tuiCmd.Flags().StringP("log", "l", "tui.log", "path for TUI debug logs")
	_ = viper.BindPFlag("client.log_file", tuiCmd.Flags().Lookup("log"))
}
