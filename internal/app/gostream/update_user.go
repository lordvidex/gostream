package gostream

import (
	"context"
	"fmt"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateUser ...
func (i *Implementation) UpdateUser(ctx context.Context, req *gostreamv1.UpdateUserRequest) (*gostreamv1.UpdateUserResponse, error) {
	user := req.User
	err := i.userRepo.UpdateUser(ctx, user)
	if err != nil {
		fmt.Println("error updating user", err)
		return nil, status.Errorf(codes.Internal, "error updating user: %v", err)
	}

	if err = i.publishUser(ctx, user); err != nil {
		fmt.Println("got error publishing data", err)
	}

	return &gostreamv1.UpdateUserResponse{}, nil
}
