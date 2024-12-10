package gostream

import (
	"context"

	"github.com/lordvidex/errs/v2/status"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// DeleteUser ...
func (i *Implementation) DeleteUser(ctx context.Context, req *gostreamv1.DeleteUserRequest) (*gostreamv1.DeleteUserResponse, error) {
	err := i.userRepo.DeleteUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Err(err)
	}

	return &gostreamv1.DeleteUserResponse{}, nil
}
