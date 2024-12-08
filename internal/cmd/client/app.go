package client

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/lordvidex/gostream/internal/config"
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// App ...
type App struct {
	cfg config.Client
}

// New ...
func New(cfg config.Client) *App {
	return &App{cfg: cfg}
}

// Watch ...
func (a *App) Watch(ctx context.Context) error {

	if a.cfg.DryRun {
		log.Println("dry run mode enabled")
		fmt.Printf("%+v\n", a.cfg)
		return nil
	}

	if a.cfg.Connections > 1 {
		return a.watchMultipleServers(ctx)
	}
	return a.watchSingleServer(ctx, a.cfg.Name)
}

func (a *App) watchMultipleServers(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(a.cfg.Connections)

	for i := range a.cfg.Connections {
		clientName := fmt.Sprintf("%s#%d", a.cfg.Name, i+1)
		g.Go(func() error {
			return a.watchSingleServer(ctx, clientName)
		})
	}

	return g.Wait()

}

func (a *App) watchSingleServer(ctx context.Context, clientName string) error {
	cl, err := a.findBestServer()
	if err != nil {
		return fmt.Errorf("couldn't find best server: %w", err)
	}

	stream, err := cl.Watch(ctx, &gostreamv1.WatchRequest{
		Identifier: getClientName(clientName),
		Entity:     a.cfg.Entities.Values(),
	})
	if err != nil {
		return fmt.Errorf("error watching: %w", err)
	}

	for {
		_, err := stream.Recv()
		fmt.Printf("client %s got message\n", clientName)
		if err != nil {
			return err
		}
	}
}

func (a *App) findBestServer() (gostreamv1.WatchersServiceClient, error) {
	// TODO: implement, for quick testing, connect to first server
	picked := rand.IntN(len(a.cfg.Servers))
	addr := a.cfg.Servers[picked]
	fmt.Println("picked server", addr)

	conn, err := grpc.NewClient(addr,
		grpc.WithBlock(),
		grpc.WithTimeout(a.cfg.ConnectionTimeout),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return gostreamv1.NewWatchersServiceClient(conn), nil
}

func getClientName(name string) *string {
	if name == "" {
		return nil
	}
	return &name
}
