package gostream

import (
	"context"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListPets ...
func (i *Implementation) ListPets(ctx context.Context, req *gostreamv1.ListPetsRequest) (*gostreamv1.ListPetsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPetsResponse stub")
}
