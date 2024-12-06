package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/urfave/cli/v3"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/lordvidex/gostream/internal/app/gostream"
	"github.com/lordvidex/gostream/internal/db/pg"
	"github.com/lordvidex/gostream/internal/watchers"
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// Serve ...
func serveFn(ctx context.Context, _ *cli.Command) error {
	repo, err := pg.NewRepository(ctx, config.DSN, pg.WithRunMigrations(config.RunMigration))
	if err != nil {
		return err
	}

	clientPub := &watchers.Client{}
	serverPub, err := watchers.NewServer(ctx, config.RedisURL, clientPub)
	if err != nil {
		return fmt.Errorf("error creating serverPubSub: %w", err)
	}

	srv := gostream.NewService(repo, clientPub, serverPub)
	s := grpc.NewServer()
	gostreamv1.RegisterPetServiceServer(s, srv)
	gostreamv1.RegisterUserServiceServer(s, srv)
	gostreamv1.RegisterWatchersServiceServer(s, srv)
	reflection.Register(s)

	addr := fmt.Sprintf(":%d", config.GRPCPort)
	log.Println("server listening on", addr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := s.Serve(lis); err != nil {
			return fmt.Errorf("server stopped with err: %w", err)
		}
		return nil
	})
	g.Go(func() error {
		if err := serveHTTPGateway(ctx, addr); err != nil {
			return fmt.Errorf("server HTTP gateway stopped with err: %w", err)
		}
		return nil
	})

	return g.Wait()
}

func serveHTTPGateway(ctx context.Context, endpoint string) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := gostreamv1.RegisterPetServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return err
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", config.HTTPPort), mux)
}
