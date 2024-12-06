package gostream

import (
	"context"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

func (i *Implementation) publishPet(ctx context.Context, p *gostreamv1.Pet) error {
	data := &gostreamv1.WatchResponse{
		Kind: gostreamv1.EventKind_EVENT_KIND_UPDATE,
		Data: &gostreamv1.WatchResponse_Update{
			Update: &gostreamv1.WatchResponse_WatchData{
				Entity: gostreamv1.Entity_ENTITY_PET,
				Data: &gostreamv1.WatchResponse_WatchData_Pet{
					Pet: p,
				},
			},
		},
	}
	return i.publishData(ctx, data)
}

func (i *Implementation) publishUser(ctx context.Context, p *gostreamv1.User) error {
	data := &gostreamv1.WatchResponse{
		Kind: gostreamv1.EventKind_EVENT_KIND_UPDATE,
		Data: &gostreamv1.WatchResponse_Update{
			Update: &gostreamv1.WatchResponse_WatchData{
				Entity: gostreamv1.Entity_ENTITY_USER,
				Data: &gostreamv1.WatchResponse_WatchData_User{
					User: p,
				},
			},
		},
	}
	return i.publishData(ctx, data)
}

func (i *Implementation) publishData(ctx context.Context, d *gostreamv1.WatchResponse) error {
	return i.serverPub.PublishToServers(ctx, d)
}
