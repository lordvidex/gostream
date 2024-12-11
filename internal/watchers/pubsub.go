package watchers

import (
	"context"
	"fmt"

	"github.com/catalystgo/catalystgo/closer"
	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
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
	closer    closer.Closer
	done      chan struct{}
}

// NewPubSub ...
func NewPubSub(ctx context.Context, redisURL string, clientPub ClientPublisher) (*PubSub, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	s := PubSub{
		cl:        redis.NewClient(opts),
		clientPub: clientPub,
		closer:    closer.New(),
		done:      make(chan struct{}),
	}
	s.closer.AddByOrder(closer.HighOrder, func() error {
		cancel()
		return nil
	})

	go s.subscribe(ctx)

	return &s, nil
}

// Close ...
func (s *PubSub) Close() error {
	// close ctx and subscriber
	s.closer.CloseAll()
	s.closer.Wait()
	<-s.done

	// close the client
	if err := s.cl.Close(); err != nil {
		return err
	}
	fmt.Println("pubsub closed")
	return nil

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
	defer func() {
		close(s.done)
	}()
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
