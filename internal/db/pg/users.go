package pg

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/lordvidex/errs/v2"

	"github.com/lordvidex/gostream/internal/db/inmemory"
	"github.com/lordvidex/gostream/internal/entity"
	"github.com/lordvidex/gostream/pkg/md5hash"

	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

var _ inmemory.DataSource[uint64, entity.User] = (*UserDataSource)(nil)

// UserDataSource implements inmemory.DataSource
type UserDataSource struct {
	*Repository
}

func NewUserDataSource(repo *Repository) *UserDataSource {
	return &UserDataSource{repo}
}

func (s *UserDataSource) Hash(ctx context.Context) (string, error) {
	q, err := md5hash.Query(entity.User{})
	if err != nil {
		return "", err
	}

	var result string
	if err = s.pool.QueryRow(ctx, q).Scan(&result); err != nil {
		return "", err
	}
	return result, nil
}

func (s *UserDataSource) FetchAll(ctx context.Context) ([]entity.User, error) {
	res, err := s.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	fin := make([]entity.User, 0, len(res))
	for _, u := range res {
		fin = append(fin, entity.User{User: u})
	}

	return fin, nil
}

func (s *UserDataSource) Fetch(ctx context.Context, id ...uint64) ([]entity.User, error) {
	query, params, err := sq.Select("id", "name", "age", "nationality").From("stream_users").
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]entity.User, 0)
	for rows.Next() {
		var x gostreamv1.User
		if err = rows.Scan(&x.Id, &x.Name, &x.Age, &x.Nationality); err != nil {
			return nil, err
		}
		res = append(res, entity.User{User: &x})
	}
	return res, nil
}

// CreateUser ...
func (r *Repository) CreateUser(ctx context.Context, p *gostreamv1.User) (uint64, error) {
	q := sq.Insert("stream_users").
		Columns("name", "age", "nationality").
		Values(p.GetName(), p.GetAge(), p.GetNationality()).
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

// ListUsers ...
func (r *Repository) ListUsers(ctx context.Context) ([]*gostreamv1.User, error) {
	q := sq.Select("id", "name", "age", "nationality").From("stream_users").
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
	result := make([]*gostreamv1.User, 0)
	for rows.Next() {
		var u gostreamv1.User
		if err = rows.Scan(&u.Id, &u.Name, &u.Age, &u.Nationality); err != nil {
			fmt.Println(err)
		}

		result = append(result, &u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, err
}

// UpdateUser ...
func (r *Repository) UpdateUser(ctx context.Context, p *gostreamv1.User) error {
	query, params, err := sq.Update("stream_users").Where(sq.Eq{"id": p.GetId()}).
		PlaceholderFormat(sq.Dollar).
		Set("name", p.GetName()).
		Set("age", p.GetAge()).
		Set("nationality", p.GetNationality()).
		Suffix("RETURNING id").ToSql()

	if err != nil {
		return errs.B().Code(errs.Internal).Msg("sq error").Err()
	}
	if err = r.pool.QueryRow(ctx, query, params...).Scan(&p.Id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.B().Code(errs.InvalidArgument).Show().Msg("user does not exist").Err()
		}
		return errs.WrapCode(err, errs.Internal, "database error occurred")
	}
	return nil
}

// DeleteUser ...
func (r *Repository) DeleteUser(ctx context.Context, p uint64) error {
	query, params, err := sq.Delete("stream_users").Where(sq.Eq{"id": p}).
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return errs.B().Code(errs.Internal).Msg("sq error").Err()
	}
	if _, err = r.pool.Exec(ctx, query, params...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.B().Code(errs.InvalidArgument).Show().Msg("user does not exist").Err()
		}
		return errs.WrapCode(err, errs.Internal, "database error occurred")
	}
	return nil
}
