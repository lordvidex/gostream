package gostream

import (
	"context"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListUsers ...
func (i *Implementation) ListUsers(ctx context.Context, req *gostreamv1.ListUsersRequest) (*gostreamv1.ListUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUsersResponse stub")
}
