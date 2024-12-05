package pg

import (
	"context"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

func (r *Repository) CreatePet(ctx context.Context, p *gostreamv1.Pet) (uint64, error) {
	// TODO: use squirrel and pgxpool to make insert queries to db
	return 0, nil
}
