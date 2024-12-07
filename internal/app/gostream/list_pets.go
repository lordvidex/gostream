package gostream

import (
	"context"
	"fmt"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListPets ...
func (i *Implementation) ListPets(ctx context.Context, req *gostreamv1.ListPetsRequest) (*gostreamv1.ListPetsResponse, error) {
	pets, err := i.petRepo.ListPets(ctx)
	if err != nil {
		fmt.Println("error listing pets", err)
		return nil, status.Errorf(codes.Internal, "error listing pets: %v", err)
	}

	return &gostreamv1.ListPetsResponse{Pets: pets}, nil
}
