package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/lordvidex/gostream"
	"github.com/pressly/goose/v3"
)

// Repository ...
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository ...
func NewRepository(ctx context.Context, connString string, opts ...ApplyOption) (*Repository, error) {
	var opt options
	for _, o := range opts {
		o(&opt)
	}

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error pinging db: %w", err)
	}

	if opt.RunMigrations {
		db := stdlib.OpenDBFromPool(pool)
		defer db.Close()
		if err = pgMigrate(db); err != nil {
			return nil, err
		}
	}

	return &Repository{pool: pool}, nil
}

type options struct {
	RunMigrations bool
}

// ApplyOption ...
type ApplyOption func(*options)

// WithRunMigrations ...
func WithRunMigrations(value bool) ApplyOption {
	return func(x *options) {
		x.RunMigrations = value
	}
}

// Close ...
func (r *Repository) Close() error {
	r.pool.Close()
	return nil
}

func pgMigrate(db *sql.DB) error {
	goose.SetBaseFS(gostream.MigrationFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	fmt.Println("migration successfully completed")
	return nil
}
