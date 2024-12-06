package pg

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// CreatePet ...
func (r *Repository) CreatePet(ctx context.Context, p *gostreamv1.Pet) (uint64, error) {
	q := sq.Insert("pets").
		Columns("kind", "name", "age").
		Values(p.GetKind(), p.GetName(), p.GetAge()).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	query, params, err := q.ToSql()
	if err != nil {
		fmt.Println("sq error: ", err)
		return 0, err
	}
	if err = r.pool.QueryRow(ctx, query, params...).Scan(&p.Id); err != nil {
		return 0, err
	}

	return p.Id, nil
}
