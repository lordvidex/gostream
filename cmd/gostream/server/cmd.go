package server

import (
	"context"
	"errors"
	"fmt"

	altsrc "github.com/urfave/cli-altsrc/v3"
	"github.com/urfave/cli/v3"

	pkgcli "github.com/lordvidex/gostream/pkg/cli"
)

var config struct {
	GRPCPort int64
	HTTPPort int64
	LogFile  string
	DryRun   bool
	DSN      string
	RedisURL string
	// CfgFilePath is the path to config.toml file if provided
	CfgFilePath  string
	RunMigration bool
}

var Cmd = &cli.Command{
	Name: "server",
	Commands: []*cli.Command{
		{
			Name: "serve",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:        "grpc",
					Usage:       "grpc port",
					Value:       50051,
					OnlyOnce:    true,
					Destination: &config.GRPCPort,
				},
				&cli.IntFlag{
					Name:        "http",
					Usage:       "http port",
					Value:       8080,
					OnlyOnce:    true,
					Destination: &config.HTTPPort,
				},
				&cli.StringFlag{
					Name:        "log",
					Aliases:     []string{"l"},
					Usage:       "path to log file",
					DefaultText: "stdout",
					Destination: &config.LogFile,
				},
				&cli.BoolFlag{
					Name:        "migrations",
					Aliases:     []string{"m"},
					Usage:       "run migrations Up before starting application",
					Destination: &config.RunMigration,
				},
				&cli.StringFlag{
					Name:     "dsn",
					Usage:    "database full URL",
					Required: true,
					Aliases:  []string{"d"},
					Sources: pkgcli.MergeChains(
						cli.NewValueSourceChain(cli.EnvVar("DSN")),
						altsrc.TOML("server.dsn", config.CfgFilePath, "config.toml"), // not working
					),
					Destination: &config.DSN,
				},
				&cli.StringFlag{
					Name:    "redis",
					Aliases: []string{"r"},
					Usage:   "redis full URL",
					Sources: pkgcli.MergeChains(
						cli.NewValueSourceChain(cli.EnvVar("REDIS_URL")),
						altsrc.TOML("server.redis_url", config.CfgFilePath, "config.toml"),
					),
					Destination: &config.RedisURL,
				},
				&cli.BoolFlag{
					Name:  "dry-run",
					Usage: "print configs and exit",
					Action: func(ctx context.Context, _ *cli.Command, v bool) error {
						if v {
							fmt.Printf("%+v\n", config)
							return errors.New("exiting")
						}
						return nil
					},
					Destination: &config.DryRun,
				},
			},
			Action: serveFn,
			Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
				fmt.Println("before called, cfgfilepath saved")
				config.CfgFilePath = cmd.String("config")
				return ctx, nil
			},
		},
	},
}
