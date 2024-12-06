package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/lordvidex/gostream/internal/app/gostream"
	"github.com/lordvidex/gostream/internal/db/pg"
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Serve ...
func serveFn(ctx context.Context, _ *cli.Command) error {
	repo, err := pg.NewRepository(ctx, config.DSN, pg.WithRunMigrations(config.RunMigration))
	if err != nil {
		return err
	}

	srv := gostream.NewService(repo, nil, nil)
	s := grpc.NewServer()
	gostreamv1.RegisterPetServiceServer(s, srv)
	gostreamv1.RegisterUserServiceServer(s, srv)
	gostreamv1.RegisterWatchersServiceServer(s, srv)
	reflection.Register(s)

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
