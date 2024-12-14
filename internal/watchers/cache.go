package watchers

import (
	"context"
	"fmt"

	"github.com/lordvidex/gostream/internal/entity"
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// Cache is a wrapper for updatePet callback and updateUsers callback that
// sends snapshot of new cache data over to listening clients when datasource is updated.
type Cache struct {
	ctx context.Context
	cl  ClientPublisher
}

func NewCache(ctx context.Context, cl ClientPublisher) *Cache {
	return &Cache{ctx, cl}
}

func (c *Cache) UpdatePets(kind gostreamv1.EventKind, data []entity.Pet) {
	switch kind {
	case gostreamv1.EventKind_EVENT_KIND_UPDATE, gostreamv1.EventKind_EVENT_KIND_DELETE:
		for _, pet := range data {
			d := &gostreamv1.WatchResponse{
				Kind:   kind,
				Entity: gostreamv1.Entity_ENTITY_PET,
				Data: &gostreamv1.WatchResponse_Update{
					Update: &gostreamv1.WatchResponse_WatchData{
						Data: &gostreamv1.WatchResponse_WatchData_Pet{
							Pet: pet.Pet,
						},
					},
				},
			}
			err := c.cl.PublishToClients(c.ctx, d)
			if err != nil {
				fmt.Println(err)
			}
		}
	case gostreamv1.EventKind_EVENT_KIND_SNAPSHOT:
		d := &gostreamv1.WatchResponse{
			Kind:   kind,
			Entity: gostreamv1.Entity_ENTITY_PET,
			Data: &gostreamv1.WatchResponse_Snapshot{
				Snapshot: &gostreamv1.WatchResponse_WatchSnapshot{
					Snapshot: descPets(data),
				},
			},
		}
		err := c.cl.PublishToClients(c.ctx, d)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func descPets(data []entity.Pet) []*gostreamv1.WatchResponse_WatchData {
	arr := make([]*gostreamv1.WatchResponse_WatchData, len(data))
	for i, v := range data {
		arr[i] = &gostreamv1.WatchResponse_WatchData{
			Data: &gostreamv1.WatchResponse_WatchData_Pet{
				Pet: v.Pet,
			},
		}
	}
	return arr
}

func (c *Cache) UpdateUsers(kind gostreamv1.EventKind, data []entity.User) {
	switch kind {
	case gostreamv1.EventKind_EVENT_KIND_UPDATE, gostreamv1.EventKind_EVENT_KIND_DELETE:
		for _, user := range data {
			d := &gostreamv1.WatchResponse{
				Kind:   kind,
				Entity: gostreamv1.Entity_ENTITY_USER,
				Data: &gostreamv1.WatchResponse_Update{
					Update: &gostreamv1.WatchResponse_WatchData{
						Data: &gostreamv1.WatchResponse_WatchData_User{
							User: user.User,
						},
					},
				},
			}
			err := c.cl.PublishToClients(c.ctx, d)
			if err != nil {
				fmt.Println(err)
			}
		}
	case gostreamv1.EventKind_EVENT_KIND_SNAPSHOT:
		d := &gostreamv1.WatchResponse{
			Kind:   kind,
			Entity: gostreamv1.Entity_ENTITY_USER,
			Data: &gostreamv1.WatchResponse_Snapshot{
				Snapshot: &gostreamv1.WatchResponse_WatchSnapshot{
					Snapshot: descUsers(data),
				},
			},
		}
		err := c.cl.PublishToClients(c.ctx, d)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func descUsers(data []entity.User) []*gostreamv1.WatchResponse_WatchData {
	arr := make([]*gostreamv1.WatchResponse_WatchData, len(data))
	for i, v := range data {
		arr[i] = &gostreamv1.WatchResponse_WatchData{
			Data: &gostreamv1.WatchResponse_WatchData_User{
				User: v.User,
			},
		}
	}
	return arr
}
