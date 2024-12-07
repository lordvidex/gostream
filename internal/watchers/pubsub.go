package watchers

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
	"google.golang.org/protobuf/proto"
)

const redisChannel = "gostream"

// ClientPublisher ...
type ClientPublisher interface {
	PublishToClients(context.Context, *gostreamv1.WatchResponse) error
}

// PubSub listens for updates from other servers
type PubSub struct {
	cl        redis.UniversalClient
	clientPub ClientPublisher
	// TODO: add some local storage
}

// NewPubSub ...
func NewPubSub(ctx context.Context, redisURL string, clientPub ClientPublisher) (*PubSub, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	s := PubSub{cl: redis.NewClient(opts), clientPub: clientPub}

	go s.subscribe(ctx)

	return &s, nil
}

// Close ...
func (s *PubSub) Close() error {
	return s.cl.Close()
}

// PublishToServers ...
func (s *PubSub) PublishToServers(ctx context.Context, data *gostreamv1.WatchResponse) error {
	b, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	return s.cl.Publish(ctx, redisChannel, string(b)).Err()
}

func (s *PubSub) subscribe(ctx context.Context) {
	sub := s.cl.Subscribe(ctx, redisChannel)
	defer sub.Close()

	for {
		select {
		// program is ending
		case <-ctx.Done():
			fmt.Println("stopping server subscribing")
			return

		case msg := <-sub.Channel():
			var data gostreamv1.WatchResponse
			if err := proto.Unmarshal([]byte(msg.Payload), &data); err != nil {
				fmt.Println("error unmarshalling proto data.. skipped", err)
				fmt.Println("text data received: ", msg.Payload)
				fmt.Println("payload slice", msg.PayloadSlice)
				break
			}

			if err := s.clientPub.PublishToClients(context.Background(), &data); err != nil {
				fmt.Println("error publishing data to clients from server subscription.. skipped", err)
				break
			}

			fmt.Println("data sent to client from server subscription.", "success")
		}
	}
}
