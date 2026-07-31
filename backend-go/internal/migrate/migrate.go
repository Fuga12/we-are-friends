package migrate

import (
	"context"
	"fmt"
	"strings"

	"github.com/Fuga12/we-are-friends/backend-go/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Up applies the initial schema exactly once. The advisory lock prevents two
// API replicas from attempting the first migration simultaneously.
func Up(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, `select pg_advisory_lock(91726301)`); err != nil {
		return err
	}
	defer pool.Exec(context.Background(), `select pg_advisory_unlock(91726301)`)
	if _, err := pool.Exec(ctx, `create table if not exists app_schema_versions (version integer primary key)`); err != nil {
		return err
	}
	var applied bool
	if err := pool.QueryRow(ctx, `select exists(select 1 from app_schema_versions where version=1)`).Scan(&applied); err != nil {
		return err
	}
	if applied {
		return nil
	}
	var mediaExists bool
	if err := pool.QueryRow(ctx, `select to_regclass('public.media') is not null`).Scan(&mediaExists); err != nil {
		return err
	}
	if !mediaExists {
		upSQL := strings.Split(migrations.Initial, "-- +goose Down")[0]
		for _, statement := range strings.Split(upSQL, ";") {
			if strings.TrimSpace(statement) == "" {
				continue
			}
			if _, err := pool.Exec(ctx, statement); err != nil {
				return fmt.Errorf("apply initial schema: %w", err)
			}
		}
	}
	_, err := pool.Exec(ctx, `insert into app_schema_versions(version) values (1) on conflict do nothing`)
	return err
}
