package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

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

	data := "{}"
	if len(event.Data) > 0 {
		b, err := json.Marshal(event.Data)
		if err != nil {
			return fmt.Errorf("unable to marshal data to json: %w", err)
		}

		data = string(b)
	}

	txn, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}

	if _, err := txn.ExecContext(ctx, stmt, event.Hash, event.ErrorId, event.Message, ts, data); err != nil {
		return fmt.Errorf("inserting error into clickhouse: %w", err)
	}

	if err := txn.Commit(); err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}

	return nil
}

func (d *Database) SelectErrorsByCount(ctx context.Context) ([]*protos.QueryResponse_Result, error) {
	const stmt = `
SELECT
	hash, 
	any(error_id) as error_id,
	any(message) as message,
	min(timestamp) as first_seen,
	max(timestamp) as last_seen,
	any(data) as data,
	count() as count 
FROM 
	errors
GROUP BY
	hash
ORDER BY
	count DESC`

	iter, err := d.db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("running query: %w", err)
	}

	results := []*protos.QueryResponse_Result{}

	for iter.Next() {
		var (
			hash      string
			errorID   string
			message   string
			firstSeen time.Time
			lastSeen  time.Time
			data      string
			count     int64
		)
		if err := iter.Scan(&hash, &errorID, &message, &firstSeen, &lastSeen, &data, &count); err != nil {
			return nil, fmt.Errorf("scanning results: %w", err)
		}

		var mapData map[string]string
		if err := json.Unmarshal([]byte(data), &mapData); err != nil {
			return nil, fmt.Errorf("unmarshalling json data to map: %w", err)
		}

		error := &protos.Error{
			Hash:    hash,
			ErrorId: errorID,
			Message: message,
			Data:    mapData,
		}

		firstSeenTS, err := ptypes.TimestampProto(firstSeen)
		if err != nil {
			return nil, fmt.Errorf("converting first seen from time.Time to *timestamp.Timestamp: %w", err)
		}

		lastSeenTS, err := ptypes.TimestampProto(lastSeen)
		if err != nil {
			return nil, fmt.Errorf("converting last seen from time.Time to *timestamp.Timestamp: %w", err)
		}

		results = append(results, &protos.QueryResponse_Result{
			Error:     error,
			FirstSeen: firstSeenTS,
			LastSeen:  lastSeenTS,
			Aggregation: &protos.QueryResponse_Result_Count{
				Count: count,
			},
		})
	}

	return results, nil
}
