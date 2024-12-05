package gostream

import gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"

// Implementation ...
type Implementation struct {
	gostreamv1.UnimplementedPetServiceServer
	gostreamv1.UnimplementedUserServiceServer
	gostreamv1.UnimplementedWatchersServiceServer
}

// NewService ...
func NewService() *Implementation {
	return &Implementation{}
}
