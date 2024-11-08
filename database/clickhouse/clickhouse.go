package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
)

func NewClickhouseClient(host, port, username, password, database string) (*sql.DB, error) {

	connStr := fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s", username, password, host, port, database)

	db, err := sql.Open("clickhouse", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open ClickHouse connection: %w", err)
	}

	if err := db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	fmt.Println("Successfully connected to ClickHouse!")
	return db, nil
}
