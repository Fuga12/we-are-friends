package migrate

import (
	"context"
	"fmt"
	"strings"

	"github.com/Fuga12/we-are-friends/backend-go/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
)

type migration struct {
	version int
	sql     string
}

// Up serially applies every embedded migration once. The table records each
// version independently, allowing production databases to move beyond v1.
func Up(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, `select pg_advisory_lock(91726301)`); err != nil {
		return err
	}
	defer pool.Exec(context.Background(), `select pg_advisory_unlock(91726301)`)
	if _, err := pool.Exec(ctx, `create table if not exists app_schema_versions (version integer primary key)`); err != nil {
		return err
	}
	for _, item := range []migration{{1, migrations.Initial}, {2, migrations.AdminAuthAndPublication}} {
		var applied bool
		if err := pool.QueryRow(ctx, `select exists(select 1 from app_schema_versions where version=$1)`, item.version).Scan(&applied); err != nil {
			return err
		}
		if applied {
			continue
		}
		upSQL := strings.Split(item.sql, "-- +goose Down")[0]
		for _, statement := range strings.Split(upSQL, ";") {
			if strings.TrimSpace(statement) == "" {
				continue
			}
			if _, err := pool.Exec(ctx, statement); err != nil {
				return fmt.Errorf("apply migration %d: %w", item.version, err)
			}
		}
		if _, err := pool.Exec(ctx, `insert into app_schema_versions(version) values ($1)`, item.version); err != nil {
			return err
		}
	}
	return nil
}
