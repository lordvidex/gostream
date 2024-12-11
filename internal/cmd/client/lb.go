package client

import (
	"context"
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/lordvidex/errs/v2"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// clientsLock allows gostream client to make good load balancing decisions by allowing clients connect ONLY one after the other.
// This is needed because `Advertise` responses were returning faster than we could start watching, therefore all clients were
// picking the same server.
var clientsLock sync.Mutex

type server struct {
	client gostreamv1.WatchersServiceClient
	addr   string
	load   float64
}

// String ...
func (s server) String() string {
	return fmt.Sprintf("%s(%d)", s.addr, int(s.load))
}

func (a *App) findBestServer(clientName string) (gostreamv1.WatchersServiceClient, error) {

	clientsLock.Lock()
	defer clientsLock.Unlock()

	var (
		servers  = make([]server, 0, len(a.cfg.Servers))
		totalErr error
	)

	for _, addr := range a.cfg.Servers {
		conn, err := a.cachedConn(addr)
		if err != nil {
			fmt.Println("error getting conn for addr:", addr)
			totalErr = errs.Wrap(totalErr, err)
			continue
		}
		ctx, _ := context.WithTimeout(a.ctx, a.cfg.ConnectionTimeout)
		cl := gostreamv1.NewWatchersServiceClient(conn)

		var res *gostreamv1.AdvertiseResponse
		res, err = cl.Advertise(ctx, &gostreamv1.AdvertiseRequest{
			Metrics: []gostreamv1.ServerMetric{gostreamv1.ServerMetric_SERVER_METRIC_STREAMS},
		})
		if err != nil {
			totalErr = errs.Wrap(totalErr, err)
			continue
		}
		servers = append(servers, server{client: cl, load: getScore(res), addr: addr})
	}

	if len(servers) == 0 {
		return nil, errs.B(totalErr).Code(errs.NotFound).Msg("no servers available").Err()
	}

	sort.Slice(servers, func(i, j int) bool { return servers[i].load < servers[j].load })
	fmt.Println("got servers: ", servers)
	fmt.Printf("client %s picked server: %s\n", clientName, servers[0].addr)
	return servers[0].client, nil
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
