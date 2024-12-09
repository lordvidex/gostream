package client

import (
	"context"
	"fmt"
	"math"
	"sort"
	"sync"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// TODO: fix and refactor
func (a *App) findBestServer(clientName string) (gostreamv1.WatchersServiceClient, error) {
	type server struct {
		client gostreamv1.WatchersServiceClient
		addr   string
		load   float64
	}

	var (
		servers = make([]server, 0, len(a.cfg.Servers))
		er      error
		ch      = make(chan server, 1)
	)

	wg := &sync.WaitGroup{}
	for _, addr := range a.cfg.Servers {
		wg.Add(1)
		go func() {
			defer wg.Done()

			conn, err := a.cachedConn(addr)
			if err != nil {
				fmt.Println("error getting connCache for addr:", addr)
				er = err
				return
			}
			cl := gostreamv1.NewWatchersServiceClient(conn)
			res, err := cl.Advertise(context.Background(), &gostreamv1.AdvertiseRequest{
				Metrics: []gostreamv1.ServerMetric{gostreamv1.ServerMetric_SERVER_METRIC_STREAMS},
			})
			if err != nil {
				er = err
				return
			}
			ch <- server{client: cl, load: getScore(res), addr: addr}
			servers = append(servers, server{client: cl, load: getScore(res), addr: addr})
		}()
	}

	go func() {
		// close when senders finish
		wg.Wait()
		close(ch)
	}()

	for res := range ch {
		servers = append(servers, res)
	}

	if len(servers) == 0 {
		return nil, fmt.Errorf("no servers picked: %w", er)
	}
	fmt.Println(servers)
	sort.Slice(servers, func(i, j int) bool { return servers[i].load < servers[j].load })
	fmt.Printf("client %s picked server: %s client: %v\n", clientName, servers[0].addr, servers[0].client)
	return servers[0].client, nil
}

func getScore(res *gostreamv1.AdvertiseResponse) float64 {
	score := math.MaxFloat64
	for _, v := range res.Response {
		if v.Metric == gostreamv1.ServerMetric_SERVER_METRIC_GOROUTINES {
			return v.Value
		}
	}
	return score
}
