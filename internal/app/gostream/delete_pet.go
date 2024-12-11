package gostream

import (
	"context"
	"fmt"

	"github.com/lordvidex/errs/v2/status"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// DeletePet ...
func (i *Implementation) DeletePet(ctx context.Context, req *gostreamv1.DeletePetRequest) (*gostreamv1.DeletePetResponse, error) {
	err := i.petRepo.DeletePet(ctx, req.PetId)
	if err != nil {
		return nil, status.Err(err)
	}

	i.petCache.Delete(req.PetId)
	if err = i.publishPetDelete(ctx, req.PetId); err != nil {
		fmt.Println("got error publishing delete data", err)
	}

	return &gostreamv1.DeletePetResponse{}, nil
}
