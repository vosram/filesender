package api

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/vosram/filesender/backend/internal/database"
)

type apiConfig struct {
	db        *database.Queries
	JWTSecret string
	platform  string
}

func New(dbConnString, jwtSecret, platform string) (*apiConfig, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbConnString)
	if err != nil {
		return nil, err
	}
	queries := database.New(conn)
	return &apiConfig{
		db:        queries,
		JWTSecret: jwtSecret,
		platform:  platform,
	}, nil
}
