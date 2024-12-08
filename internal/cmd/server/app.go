package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/lordvidex/gostream/internal/config"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/catalystgo/catalystgo/closer"
	"github.com/lordvidex/gostream/internal/app/gostream"
	"github.com/lordvidex/gostream/internal/db/pg"
	"github.com/lordvidex/gostream/internal/watchers"
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// App ...
type App struct {
	closer closer.Closer
	cfg    config.Server
}

// New ...
func New(cfg config.Server) *App {
	return &App{
		cfg: cfg,
		closer: closer.New(
			closer.WithSignals(os.Kill, os.Interrupt),
			closer.WithTimeout(time.Minute*3),
		),
	}
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

	a.closer.Add(clientWatcher.Close, serverPub.Close, repo.Close)

	srv := gostream.NewService(repo, serverPub, clientWatcher)

	s := grpc.NewServer()
	gostreamv1.RegisterPetServiceServer(s, srv)
	gostreamv1.RegisterUserServiceServer(s, srv)
	gostreamv1.RegisterWatchersServiceServer(s, srv)
	reflection.Register(s)

	a.closer.AddByOrder(closer.HighOrder, func() error {
		s.GracefulStop()
		fmt.Println("grpc server gracefully stopped")
		return nil
	})

	addr := fmt.Sprintf(":%d", a.cfg.GRPCPort)
	log.Println("server listening on", addr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	servers, ctx := errgroup.WithContext(ctx)
	servers.Go(func() error {
		if err = s.Serve(lis); err != nil {
			return fmt.Errorf("server stopped with err: %w", err)
		}
		return nil
	})
	servers.Go(func() error {
		if err = a.serveHTTPGateway(ctx, addr); err != nil {
			return fmt.Errorf("server HTTP gateway stopped with err: %w", err)
		}
		return nil
	})

	err = servers.Wait() // servers finish first
	a.closer.Wait()      // wait for connections to close

	return err
}

func (a *App) serveHTTPGateway(ctx context.Context, endpoint string) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := gostreamv1.RegisterPetServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return err
	}

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.cfg.HTTPPort),
		Handler: mux,
	}

	a.closer.AddByOrder(closer.HighOrder, func() error {
		shutCtx, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()
		return s.Shutdown(shutCtx)
	})

	return s.ListenAndServe()
}
