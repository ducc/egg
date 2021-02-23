package database

import (
	"context"
	"database/sql"
	"fmt"

	// clickhouse database/sql driver
	_ "github.com/ClickHouse/clickhouse-go"

	"github.com/ducc/egg/env"
	"github.com/ducc/egg/protos"
	"github.com/golang/protobuf/ptypes"
)

type Database struct {
	db *sql.DB
}

func New(ctx context.Context) (*Database, error) {
	db, err := sql.Open("clickhouse", env.ClickHouseURI())
	if err != nil {
		return nil, fmt.Errorf("connecting to clickhouse: %w", err)
	}

	if _, err := db.ExecContext(ctx, schema); err != nil {
		return nil, fmt.Errorf("ensuring schema exists: %w", err)
	}

	return &Database{
		db: db,
	}, nil
}

// InsertError takes an error and inserts it into clickhouse
func (d *Database) InsertError(ctx context.Context, event *protos.Error) error {
	const stmt = `
INSERT INTO errors (hash, error_id, message, timestamp, data) 
VALUES ($1, $2, $3, $4, $5);
`

	ts, err := ptypes.Timestamp(event.Timestamp)
	if err != nil {
		return fmt.Errorf("converting timestamp to time.Time: %w", err)
	}

	ts = ts.UTC()

	if _, err := d.db.ExecContext(ctx, stmt, event.Hash, event.ErrorId, event.Message, ts, event.Data); err != nil {
		return fmt.Errorf("inserting error into clickhouse: %w", err)
	}

	return nil
}
