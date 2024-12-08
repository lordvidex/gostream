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

// ListPets ...
func (r *Repository) ListPets(ctx context.Context) ([]*gostreamv1.Pet, error) {
	q := sq.Select("id", "kind", "name", "age").From("pets").
		PlaceholderFormat(sq.Dollar)

	query, params, err := q.ToSql()
	if err != nil {
		fmt.Println("sq error: ", err)
	}
	rows, err := r.pool.Query(ctx, query, params...)
	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()
	result := make([]*gostreamv1.Pet, 0)
	for rows.Next() {
		var u gostreamv1.Pet
		if err = rows.Scan(&u.Id, &u.Kind, &u.Name, &u.Age); err != nil {
			fmt.Println(err)
		}

		result = append(result, &u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, err
}

// UpdatePet ...
func (r *Repository) UpdatePet(ctx context.Context, p *gostreamv1.Pet) error {
	id := p.GetId()
	q := sq.Update("pets").Where("id = ?", id).
		PlaceholderFormat(sq.Dollar).
		Set("kind", p.GetKind()).
		Set("name", p.GetName()).
		Set("age", p.GetAge())

	query, params, err := q.ToSql()
	if err != nil {
		fmt.Println("sq error: ", err)
		return err
	}
	if _, err := r.pool.Exec(ctx, query, params...); err != nil {
		fmt.Println("Exec error: ", err)
		return err
	}
	return nil
}
