package gostream

import (
	"context"
	"fmt"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdatePet ...
func (i *Implementation) UpdatePet(ctx context.Context, req *gostreamv1.UpdatePetRequest) (*gostreamv1.UpdatePetResponse, error) {
	req.NewData.Id = req.PetId
	err := i.petRepo.UpdatePet(ctx, req.NewData)
	if err != nil {
		fmt.Println("error updating pet", err)
		return nil, status.Errorf(codes.Internal, "error updating pet: %v", err)
	}

	if err = i.publishPet(ctx, req.NewData); err != nil {
		fmt.Println("got error publishing data", err)
	}

	return &gostreamv1.UpdatePetResponse{}, nil
}
