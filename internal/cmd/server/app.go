package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/lordvidex/gostream/internal/config"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/lordvidex/gostream/internal/app/gostream"
	"github.com/lordvidex/gostream/internal/db/pg"
	"github.com/lordvidex/gostream/internal/watchers"
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// App ...
type App struct {
	cfg config.Server
}

// New ...
func New(cfg config.Server) *App {
	return &App{cfg: cfg}
}

// Serve ...
func (a *App) Serve(ctx context.Context) error {

	if a.cfg.DryRun {
		log.Println("dry run mode enabled")
		fmt.Printf("%+v\n", a.cfg)
		return nil
	}

	repo, err := pg.NewRepository(ctx, a.cfg.DSN, pg.WithRunMigrations(a.cfg.RunMigrations))
	if err != nil {
		return err
	}

	clientWatcher := watchers.NewWatcherRegistrar()
	serverPub, err := watchers.NewPubSub(ctx, a.cfg.RedisURL, clientWatcher)
	if err != nil {
		return fmt.Errorf("error creating serverPubSub: %w", err)
	}

	srv := gostream.NewService(repo, serverPub, clientWatcher)

	s := grpc.NewServer()
	gostreamv1.RegisterPetServiceServer(s, srv)
	gostreamv1.RegisterUserServiceServer(s, srv)
	gostreamv1.RegisterWatchersServiceServer(s, srv)
	reflection.Register(s)

	addr := fmt.Sprintf(":%d", a.cfg.GRPCPort)
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
		if err := serveHTTPGateway(ctx, a.cfg, addr); err != nil {
			return fmt.Errorf("server HTTP gateway stopped with err: %w", err)
		}
		return nil
	})

	return g.Wait()
}

func serveHTTPGateway(ctx context.Context, cfg config.Server, endpoint string) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := gostreamv1.RegisterPetServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return err
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.HTTPPort), mux)
}
