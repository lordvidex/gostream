/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "gostream client commands",
}

func init() {
	rootCmd.AddCommand(clientCmd)
}

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "connects to any streaming server and listens for events",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO:
		fmt.Printf("%+v\n", cfg.Client)
	},
}

func init() {
	clientCmd.AddCommand(watchCmd)

	watchCmd.Flags().Int("connections", 1, "number of parallel connections to make with the server. This is strictly for testing purpose, ideally only one connection should be made to server.")
	watchCmd.Flags().StringSlice("servers", nil, "list of servers that can be connected to, in the cluster")
	watchCmd.Flags().StringSliceP("entities", "e", []string{"all"}, "entities to watch separated by comma. Valid values include: `all`, `pet`, `user`")
	watchCmd.Flags().Bool("dry-run", false, "print configs and exit")

	viper.BindPFlag("client.connections", watchCmd.Flags().Lookup("connections"))
	viper.BindPFlag("client.servers", watchCmd.Flags().Lookup("servers"))
	viper.BindPFlag("client.entities", watchCmd.Flags().Lookup("entities"))
	viper.BindPFlag("client.dry_run", watchCmd.Flags().Lookup("dry-run"))
}
