package gostream

import (
	"context"
	"fmt"

	"github.com/lordvidex/errs/v2/status"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// UpdatePet ...
func (i *Implementation) UpdatePet(ctx context.Context, req *gostreamv1.UpdatePetRequest) (*gostreamv1.UpdatePetResponse, error) {
	req.NewData.Id = req.PetId
	err := i.petRepo.UpdatePet(ctx, req.NewData)
	if err != nil {
		return nil, status.Err(err)
	}

	if err = i.publishPetUpdate(ctx, req.NewData); err != nil {
		fmt.Println("got error publishing data", err)
	}

	return &gostreamv1.UpdatePetResponse{}, nil
}
