package gostream

import (
	"context"
	"fmt"

	"github.com/lordvidex/errs/v2/status"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// DeleteUser ...
func (i *Implementation) DeleteUser(ctx context.Context, req *gostreamv1.DeleteUserRequest) (*gostreamv1.DeleteUserResponse, error) {
	err := i.userRepo.DeleteUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Err(err)
	}

	if err = i.publishUserDelete(ctx, req.UserId); err != nil {
		fmt.Println("got error publishing delete data", err)
	}

	return &gostreamv1.DeleteUserResponse{}, nil
}
