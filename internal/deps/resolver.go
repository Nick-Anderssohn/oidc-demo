package deps

import (
	"context"

	"github.com/Nick-Anderssohn/oidc-demo/internal/sqlc/dal"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Resolver struct {
	DBPool  *pgxpool.Pool
	Queries *dal.Queries
}

func InitDepsResolver(ctx context.Context) (Resolver, error) {
	// assumes you are running via docker-compose
	// obviously you would want this to be configurable in real life.
	dbPool, err := pgxpool.New(ctx, "postgres://demo:demo@db:5432/demo")
	if err != nil {
		return Resolver{}, err
	}

	return Resolver{
		DBPool:  dbPool,
		Queries: dal.New(dbPool),
	}, nil
}

func (r *Resolver) Close() {
	if r.DBPool != nil {
		r.DBPool.Close()
	}
}
