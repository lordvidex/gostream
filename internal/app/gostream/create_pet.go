package gostream

import (
	"context"
	"fmt"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreatePet ...
func (i *Implementation) CreatePet(ctx context.Context, req *gostreamv1.CreatePetRequest) (*gostreamv1.CreatePetResponse, error) {
	pet := req.Pet
	id, err := i.petRepo.CreatePet(ctx, pet)
	if err != nil {
		fmt.Println("error creating pet", err)
		return nil, status.Errorf(codes.Internal, "error creating pet: %v", err)
	}

	pet.Id = id
	if err = i.publishPet(ctx, pet); err != nil {
		fmt.Println("got error publishing data", err)
	}

	return &gostreamv1.CreatePetResponse{Id: id}, nil
}
