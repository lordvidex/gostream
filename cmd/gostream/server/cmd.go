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
	Port     int64
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
					Name:        "port",
					Aliases:     []string{"p"},
					Value:       50051,
					OnlyOnce:    true,
					Destination: &config.Port,
				},
				&cli.StringFlag{
					Name:        "log-file",
					Aliases:     []string{"l"},
					Usage:       "log file path",
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
					Name:    "redis-url",
					Aliases: []string{"r"},
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
