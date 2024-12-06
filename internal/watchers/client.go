package watchers

import (
	"context"
	"fmt"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// Client stores streams for clients and updates them as needed.
type Client struct {
}

func (c *Client) PublishToClients(ctx context.Context, data *gostreamv1.WatchResponse) error {
	// TODO:
	fmt.Println("received data published to client")
	return nil
}
