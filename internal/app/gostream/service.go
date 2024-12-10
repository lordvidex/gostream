package gostream

import (
	"context"

	"github.com/lordvidex/gostream/internal/entity"
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// ServerPublisher ...
type ServerPublisher interface {
	PublishToServers(context.Context, *gostreamv1.WatchResponse) error
}

// WatchRegistrar registers watchers to receive updates ...
type WatchRegistrar interface {
	RegisterWatcher(*entity.Watcher) error
	DeleteWatcher(*entity.Watcher) error
	Count() int64
}

// Repository ...
type Repository interface {
	PetRepository
	UserRepository
}

// PetRepository ...
type PetRepository interface {
	CreatePet(context.Context, *gostreamv1.Pet) (uint64, error)
	ListPets(context.Context) ([]*gostreamv1.Pet, error)
	UpdatePet(context.Context, *gostreamv1.Pet) error
	DeletePet(context.Context, uint64) error
}

// UserRepository ...
type UserRepository interface {
	CreateUser(context.Context, *gostreamv1.User) (uint64, error)
	ListUsers(context.Context) ([]*gostreamv1.User, error)
	UpdateUser(context.Context, *gostreamv1.User) error
	DeleteUser(context.Context, uint64) error
}

// Implementation ...
type Implementation struct {
	// services
	gostreamv1.UnimplementedPetServiceServer
	gostreamv1.UnimplementedUserServiceServer
	gostreamv1.UnimplementedWatchersServiceServer
	//repos
	petRepo  PetRepository
	userRepo UserRepository
	// pubs
	serverPub ServerPublisher
	watchers  WatchRegistrar
}

// NewService ...
func NewService(r Repository, sp ServerPublisher, wr WatchRegistrar) *Implementation {
	return &Implementation{petRepo: r, userRepo: r, serverPub: sp, watchers: wr}
}
