/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/lordvidex/gostream/internal/cmd/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server commands",
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start a gostream server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := server.New(cfg.Server).Serve(cmd.Context()); err != nil {
			fmt.Println("finished with error", err)
		}
		fmt.Println("server finished")
	},
}

func init() {
	serverCmd.AddCommand(serveCmd)

	serveCmd.Flags().Int("grpc", 50051, "grpc port")
	serveCmd.Flags().Int("http", 8080, "http port")
	serveCmd.Flags().StringP("log", "l", "", "path to log file")
	serveCmd.Flags().BoolP("migrations", "m", false, "run migrations Up before starting application")
	serveCmd.Flags().StringP("dsn", "d", "", "database connection string")
	serveCmd.Flags().StringP("redis", "r", "", "redis connection string")
	serveCmd.Flags().Bool("dry-run", false, "print configs and exit")

	_ = viper.BindPFlag("server.grpc_port", serveCmd.Flags().Lookup("grpc"))
	_ = viper.BindPFlag("server.http_port", serveCmd.Flags().Lookup("http"))
	_ = viper.BindPFlag("server.dsn", serveCmd.Flags().Lookup("dsn"))
	_ = viper.BindPFlag("server.redis_url", serveCmd.Flags().Lookup("redis"))
	_ = viper.BindPFlag("server.log_file", serveCmd.Flags().Lookup("log"))
	_ = viper.BindPFlag("server.run_migrations", serveCmd.Flags().Lookup("migrations"))
	_ = viper.BindPFlag("server.dry_run", serveCmd.Flags().Lookup("dry-run"))
}
