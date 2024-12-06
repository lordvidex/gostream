package gostream

import (
	"context"
	"fmt"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListUsers ...
func (i *Implementation) ListUsers(ctx context.Context, req *gostreamv1.ListUsersRequest) (*gostreamv1.ListUsersResponse, error) {
	users, err := i.userRepo.ListUsers(ctx)
	if err != nil {
		fmt.Println("error listing users", err)
		return nil, status.Errorf(codes.Internal, "error listing users: %v", err)
	}

	return &gostreamv1.ListUsersResponse{Users: users}, nil
}
