package gostream

import (
	"context"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// ServerPublisher ...
type ServerPublisher interface {
	PublishToServers(context.Context, *gostreamv1.WatchResponse) error
}

// Repository ...
type Repository interface {
	PetRepository
	UserRepository
}

// PetRepository ...
type PetRepository interface {
	CreatePet(context.Context, *gostreamv1.Pet) (uint64, error)
}

// UserRepository ...
type UserRepository interface {
	// TODO: complete
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
}

// NewService ...
func NewService(r Repository, sp ServerPublisher) *Implementation {
	return &Implementation{petRepo: r, userRepo: r, serverPub: sp}
}
