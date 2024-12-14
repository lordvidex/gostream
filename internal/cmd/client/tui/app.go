package tui

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"

	"github.com/catalystgo/catalystgo/closer"
	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/lordvidex/errs/v2"

	"github.com/lordvidex/gostream/internal/config"
	"github.com/lordvidex/gostream/internal/entity"
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

type App struct {
	*Client
	ctx context.Context
	cfg config.Client
	p   *tea.Program
	// cache  *inmemory.Cache[uint64, entity.Pet]
	closer closer.Closer
}

func New(ctx context.Context, cfg config.Client) *App {
	ap := App{
		ctx:    ctx,
		Client: &Client{ctx: ctx},
		cfg:    cfg,
		closer: closer.New(),
	}

	// no datasource, so no errors
	// ap.cache, _ = inmemory.NewCache(ctx, inmemory.NewArray[uint64, entity.Pet]())
	return &ap
}

func (a *App) Run() error {
	if a.cfg.LogFile != "" {
		f, err := tea.LogToFile(a.cfg.LogFile, "debug")
		if err != nil {
			return fmt.Errorf("error logging to file: %w", err)
		}
		defer f.Close()
	}

	a.p = tea.NewProgram(newHome(a.Client))

	watchCl, err := a.connectBestServer()
	if err != nil {
		return err
	}

	go a.stream(watchCl)

	if _, err = a.p.Run(); err != nil {
		return err
	}
	return nil
}

func (a *App) stream(cl gostreamv1.WatchersServiceClient) {
	for {
		stream, err := cl.Watch(a.ctx, &gostreamv1.WatchRequest{
			Entity: []gostreamv1.Entity{
				gostreamv1.Entity_ENTITY_UNSPECIFIED,
			},
		})
		if err != nil {
			log.Println("error connecting to server")
			os.Exit(1)
		}

		for {
			v, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					log.Println("stream finished")
					return
				}
				log.Println("got error: ", err)
				break
			}

			b, err := protojson.Marshal(v)
			if err != nil {
				log.Println("error marshalling message:", err)
			}

			msgs := []tea.Msg{logMsg{Source: "server", JSON: string(b)}}

			if v.GetEntity() == gostreamv1.Entity_ENTITY_PET {
				switch v.Kind {
				case gostreamv1.EventKind_EVENT_KIND_DELETE:
					upd := v.GetUpdate()
					msgs = append(msgs, deleteMsg{ID: upd.GetPet().Id})
				case gostreamv1.EventKind_EVENT_KIND_SNAPSHOT:
					snap := v.GetSnapshot().GetSnapshot()
					res := make([]entity.Pet, 0, len(snap))
					for _, pet := range snap {
						res = append(res, entity.Pet{Pet: pet.GetPet()})
					}
					msgs = append(msgs, snapshotMsg{Pets: res})
				case gostreamv1.EventKind_EVENT_KIND_UPDATE:
					upd := v.GetUpdate()
					msgs = append(msgs, updateMsg{Pet: entity.Pet{Pet: upd.GetPet()}})
				}
			}

			for _, msg := range msgs {
				a.p.Send(msg)
			}
		}
	}

}

func (a *App) connectBestServer() (gostreamv1.WatchersServiceClient, error) {
	servers := make([]server, 0, len(a.cfg.Servers))
	var totalErr error
	for _, addr := range a.cfg.Servers {
		srv, err := a.connectServer(addr)
		if err != nil {
			totalErr = errs.Wrap(totalErr, err)
			continue
		}
		servers = append(servers, *srv)
	}

	if len(servers) == 0 {
		return nil, errs.B(totalErr).Code(errs.NotFound).Msg("no servers available").Err()
	}

	sort.Slice(servers, func(i, j int) bool { return servers[i].load < servers[j].load })
	a.petCl = servers[0].petCl
	return servers[0].watchCl, nil
}

func (a *App) connectServer(addr string) (*server, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithBlock(),
		grpc.WithTimeout(a.cfg.ConnectionTimeout),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	a.closer.Add(conn.Close)

	ctx, cancel := context.WithTimeout(a.ctx, a.cfg.ConnectionTimeout)
	defer cancel()

	cl := gostreamv1.NewWatchersServiceClient(conn)
	res, err := cl.Advertise(ctx, &gostreamv1.AdvertiseRequest{Metrics: []gostreamv1.ServerMetric{gostreamv1.ServerMetric_SERVER_METRIC_STREAMS}})
	if err != nil {
		return nil, err
	}
	return &server{watchCl: cl, petCl: gostreamv1.NewPetServiceClient(conn), load: getScore(res), addr: addr}, nil
}

type server struct {
	watchCl gostreamv1.WatchersServiceClient
	petCl   gostreamv1.PetServiceClient
	addr    string
	load    float64
}

func getScore(res *gostreamv1.AdvertiseResponse) float64 {
	score := math.MaxFloat64
	for _, v := range res.Response {
		if v.Metric == gostreamv1.ServerMetric_SERVER_METRIC_STREAMS {
			return v.Value
		}
	}
	return score
}
