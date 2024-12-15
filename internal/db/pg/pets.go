package pg

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/lordvidex/errs/v2"

	"github.com/lordvidex/gostream/pkg/md5hash"

	"github.com/lordvidex/gostream/internal/db/inmemory"
	"github.com/lordvidex/gostream/internal/entity"
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

var _ inmemory.DataSource[uint64, entity.Pet] = (*PetDataSource)(nil)

// PetDataSource implements inmemory.DataSource
type PetDataSource struct {
	*Repository
}

// NewPetDataSource ...
func NewPetDataSource(repo *Repository) *PetDataSource {
	return &PetDataSource{repo}
}

func (p *PetDataSource) Hash(ctx context.Context) (string, error) {
	q, err := md5hash.Query(entity.Pet{})
	if err != nil {
		return "", err
	}
	var result string
	if err = p.pool.QueryRow(ctx, q).Scan(&result); err != nil {
		return "", err
	}
	return result, nil
}

func (p *PetDataSource) FetchAll(ctx context.Context) ([]entity.Pet, error) {
	res, err := p.ListPets(ctx)
	if err != nil {
		return nil, err
	}

	fin := make([]entity.Pet, 0, len(res))
	for _, pet := range res {
		fin = append(fin, entity.Pet{Pet: pet})
	}
	return fin, nil
}

func (p *PetDataSource) Fetch(ctx context.Context, id ...uint64) ([]entity.Pet, error) {
	query, params, err := sq.Select("id", "kind", "name", "age").
		From("pets").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := p.pool.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]entity.Pet, 0)
	for rows.Next() {
		var x gostreamv1.Pet
		if err = rows.Scan(&x.Id, &x.Kind, &x.Name, &x.Age); err != nil {
			return nil, err
		}
		res = append(res, entity.Pet{Pet: &x})
	}
	return res, nil
}

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
	q := sq.Update("pets").Where("id = ?", p.GetId()).
		PlaceholderFormat(sq.Dollar).
		Set("kind", p.GetKind()).
		Set("name", p.GetName()).
		Set("age", p.GetAge()).Suffix("RETURNING id")

	query, params, err := q.ToSql()
	if err != nil {
		return errs.B().Code(errs.Internal).Msg("sq error").Err()
	}
	if err = r.pool.QueryRow(ctx, query, params...).Scan(&p.Id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.B().Code(errs.InvalidArgument).Msg("pet does not exist").Err()
		}
		return err
	}
	return nil
}

// DeletePet ...
func (r *Repository) DeletePet(ctx context.Context, p uint64) error {
	query, params, err := sq.Delete("pets").Where("id = ?", p).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return errs.B().Code(errs.Internal).Msg("sq error").Err()
	}
	if _, err = r.pool.Exec(ctx, query, params...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.B().Code(errs.InvalidArgument).Msg("pet does not exist").Err()
		}
		return errs.WrapCode(err, errs.Internal, "database error occurred")
	}
	return nil
}
