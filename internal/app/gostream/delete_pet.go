package gostream

import (
	"context"

	"github.com/lordvidex/errs/v2/status"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// DeletePet ...
func (i *Implementation) DeletePet(ctx context.Context, req *gostreamv1.DeletePetRequest) (*gostreamv1.DeletePetResponse, error) {
	err := i.petRepo.DeletePet(ctx, req.PetId)
	if err != nil {
		return nil, status.Err(err)
	}

	return &gostreamv1.DeletePetResponse{}, nil
}
