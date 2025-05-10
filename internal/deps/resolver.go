package deps

import (
	"context"

	"github.com/Nick-Anderssohn/oidc-demo/internal/config"
	"github.com/Nick-Anderssohn/oidc-demo/internal/sqlc/dal"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Resolver struct {
	DBPool  *pgxpool.Pool
	Queries *dal.Queries
	Config  *config.Config
}

func InitDepsResolver(ctx context.Context) (Resolver, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return Resolver{}, err
	}

	dbPool, err := pgxpool.New(ctx, cfg.PostgresConfig.ConnectionString())
	if err != nil {
		return Resolver{}, err
	}

	queries := dal.New(dbPool)

	return Resolver{
		DBPool:  dbPool,
		Queries: queries,
		Config:  &cfg,
	}, nil
}

func (r *Resolver) Close() {
	if r.DBPool != nil {
		r.DBPool.Close()
	}
}
