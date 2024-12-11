package gostream

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lordvidex/gostream/internal/entity"
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// Watch ...
func (i *Implementation) Watch(req *gostreamv1.WatchRequest, stream gostreamv1.WatchersService_WatchServer) error {
	updates := make(chan *gostreamv1.WatchResponse, 1)
	watcher := entity.NewWatcher(updates, req.GetEntity(), entity.WithWatcherIdentifier(req.GetIdentifier()))

ENT:
	for _, ent := range watcher.Entities() {
		switch ent {
		case gostreamv1.Entity_ENTITY_UNSPECIFIED:
			userSnap := i.userCache.Snapshot()
			petSnap := i.petCache.Snapshot()
			if err := stream.Send(makePetSnapshot(petSnap)); err != nil {
				return status.Error(codes.Internal, err.Error())
			}
			if err := stream.Send(makeUserSnapshot(userSnap)); err != nil {
				return status.Error(codes.Internal, err.Error())
			}
			break ENT // send snapshot once for wildcard subscribers
		case gostreamv1.Entity_ENTITY_USER:
			userSnap := i.userCache.Snapshot()
			if err := stream.Send(makeUserSnapshot(userSnap)); err != nil {
				return status.Error(codes.Internal, err.Error())
			}
		case gostreamv1.Entity_ENTITY_PET:
			petSnap := i.petCache.Snapshot()
			if err := stream.Send(makePetSnapshot(petSnap)); err != nil {
				return status.Error(codes.Internal, err.Error())
			}
		}
	}

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
