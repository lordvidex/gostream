package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"

	"github.com/lordvidex/gostream/internal/app/gostream"
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

var config struct {
	Port    int64
	LogFile string
	DryRun  bool
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
