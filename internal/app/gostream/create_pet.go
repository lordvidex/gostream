package gostream

import (
	"context"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreatePet ...
func (i *Implementation) CreatePet(ctx context.Context, req *gostreamv1.CreatePetRequest) (*gostreamv1.CreatePetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePet stub")
}
