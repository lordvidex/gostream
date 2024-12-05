package gostream

import (
	"context"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateUser ...
func (i *Implementation) CreateUser(ctx context.Context, req *gostreamv1.CreateUserRequest) (*gostreamv1.CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser stub")
}
