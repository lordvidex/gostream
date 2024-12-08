package gostream

import (
	"fmt"

	"github.com/lordvidex/gostream/internal/entity"
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Watch ...
func (i *Implementation) Watch(req *gostreamv1.WatchRequest, stream gostreamv1.WatchersService_WatchServer) error {
	updates := make(chan *gostreamv1.WatchResponse, 1)
	watcher := entity.NewWatcher(updates, req.GetEntity())

	if err := i.watchers.RegisterWatcher(watcher); err != nil {
		return status.Errorf(codes.Internal, "error registering watcher: %v", err)
	}

	defer func() {
		if err := i.watchers.DeleteWatcher(watcher); err != nil {
			fmt.Println("error deleting watcher", err)
		}
	}()

	for {
		select {
		case update, ok := <-updates:
			if !ok {
				return status.Error(codes.Canceled, "updates channel closed")
			}
			if err := stream.Send(update); err != nil {
				return status.Error(codes.Canceled, "error sending update to stream")
			}
		case <-stream.Context().Done():
			return status.Error(codes.Canceled, "stream canceled")
		}
	}
}
