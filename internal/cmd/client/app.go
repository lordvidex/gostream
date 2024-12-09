package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/catalystgo/catalystgo/closer"

	"github.com/lordvidex/gostream/internal/config"
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// App ...
type App struct {
	closer    closer.Closer
	connCache sync.Map
	cfg       config.Client
}

// New ...
func New(cfg config.Client) *App {
	return &App{
		cfg:    cfg,
		closer: closer.New(closer.WithSignals(os.Kill, os.Interrupt)),
	}
}

// Watch ...
func (a *App) Watch(ctx context.Context) error {

	if a.cfg.DryRun {
		log.Println("dry run mode enabled")
		fmt.Printf("%+v\n", a.cfg)
		return nil
	}

	a.closer.Add(func() error {
		a.connCache.Range(func(_, value any) bool {
			if err := value.(*grpc.ClientConn).Close(); err != nil {
				fmt.Println("failed tearing down *grpc.ClientConn", err)
			}
			return true
		})
		a.connCache.Clear()
		return nil
	})

	var err error
	if a.cfg.Connections > 1 {
		err = a.watchMultipleServers(ctx)
	} else {
		err = a.watchSingleServer(ctx, a.cfg.Name)
	}

	a.closer.CloseAll()
	a.closer.Wait()
	return err
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
	cl, err := a.findBestServer(clientName)
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
		v, err := stream.Recv()
		b, _ := protojson.Marshal(v)
		fmt.Printf("client %s got message: %v\n", clientName, string(b))

		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("stream finished")
				return nil
			}
			return err
		}
	}
}

func (a *App) cachedConn(addr string) (*grpc.ClientConn, error) {
	if conn, ok := a.connCache.Load(addr); ok {
		return conn.(*grpc.ClientConn), nil
	}

	conn, err := grpc.NewClient(addr,
		grpc.WithBlock(),
		grpc.WithTimeout(a.cfg.ConnectionTimeout),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	a.connCache.Store(addr, conn)
	return conn, nil
}

func getClientName(name string) *string {
	if name == "" {
		return nil
	}
	return &name
}
