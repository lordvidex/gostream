package gostream

import (
	"context"
	"fmt"

	"github.com/lordvidex/errs/v2/status"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
	"google.golang.org/grpc/codes"
)

// ListPets ...
func (i *Implementation) ListPets(ctx context.Context, req *gostreamv1.ListPetsRequest) (*gostreamv1.ListPetsResponse, error) {
	if req.Cached {
		pets := i.petCache.Snapshot()
		res := make([]*gostreamv1.Pet, 0, len(pets))
		for _, pet := range pets {
			res = append(res, pet.Pet)
		}
		return &gostreamv1.ListPetsResponse{Pets: res}, nil
	}

	pets, err := i.petRepo.ListPets(ctx)
	if err != nil {
		fmt.Println("error listing pets", err)
		return nil, status.Errorf(codes.Internal, "error listing pets: %v", err)
	}

	return &gostreamv1.ListPetsResponse{Pets: pets}, nil
}
