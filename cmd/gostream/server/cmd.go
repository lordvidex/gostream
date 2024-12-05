package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	altsrc "github.com/urfave/cli-altsrc/v3"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"

	pkgcli "github.com/lordvidex/gostream/pkg/cli"

	"github.com/lordvidex/gostream/internal/app/gostream"
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

var config struct {
	Port     int64
	LogFile  string
	DryRun   bool
	DSN      string
	RedisURL string
	// CfgFilePath is the path to config.toml file if provided
	CfgFilePath string
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

// Serve ...
func serveFn(ctx context.Context, _ *cli.Command) error {
	srv := gostream.NewService()
	s := grpc.NewServer()
	gostreamv1.RegisterPetServiceServer(s, srv)
	gostreamv1.RegisterUserServiceServer(s, srv)
	gostreamv1.RegisterWatchersServiceServer(s, srv)

	addr := fmt.Sprintf(":%d", config.Port)
	log.Println("server listening on", addr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("server stopped with err: %w", err)
	}

	return nil
}
