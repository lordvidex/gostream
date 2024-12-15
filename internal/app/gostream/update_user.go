package gostream

import (
	"context"
	"fmt"

	"github.com/lordvidex/errs/v2/status"

	"github.com/lordvidex/gostream/internal/entity"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// UpdateUser ...
func (i *Implementation) UpdateUser(ctx context.Context, req *gostreamv1.UpdateUserRequest) (*gostreamv1.UpdateUserResponse, error) {
	req.NewData.Id = req.UserId
	err := i.userRepo.UpdateUser(ctx, req.NewData)
	if err != nil {
		return nil, status.Err(err)
	}

	i.userCache.Store(entity.User{User: req.NewData})
	if err = i.publishUserUpdate(ctx, req.NewData); err != nil {
		fmt.Println("got error publishing data", err)
	}

	return &gostreamv1.UpdateUserResponse{}, nil
}
