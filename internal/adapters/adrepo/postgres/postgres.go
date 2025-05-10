package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Добавьте этот импорт
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

type PgConfig struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     string `env:"PORT" envDefault:"5432"`
	Database string `env:"DATABASE" envDefault:"postgres"`
	Username string `env:"USERNAME" envDefault:"mac"`
	Password string `env:"PASSWORD" envDefault:"1234"`
}

func New(ctx context.Context, cfg PgConfig) (*pgx.Conn, error) {
	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	migrationPath := "file://internal/adapters/adrepo/migrations"
	m, err := migrate.New(migrationPath, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to create migrations: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("unable to run migrations: %w", err)
	}
	return conn, nil
}
