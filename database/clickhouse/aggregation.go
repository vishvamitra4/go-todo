package clickhouse

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

func AggregateMetricsFromClickHouse() (interface{}, error) {

	client, err := NewClickhouseClient(os.Getenv("CH_HOST"), os.Getenv("CH_PORT"), os.Getenv("CH_USERNAME"), os.Getenv("CH_PASSWORD"), os.Getenv("CH_DB"))
	if err != nil {
		return nil, fmt.Errorf("ClickHouse client is not initialized: %v", err)
	}

	query := `
        SELECT 
            status, 
            toDate(created_at) AS created_date, 
            count(*) AS totalCount,
            sum(effort_hours) AS total_effort
        FROM todos
        GROUP BY status, created_date
        ORDER BY created_date DESC, status
    `

	// Execute the query
	rows, err := client.QueryContext(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("Error executing ClickHouse query: %v", err)
	}
	defer rows.Close()

	var result []bson.M

	for rows.Next() {
		var status string
		var createdDate string
		var totalCount int64
		var effortHours int64

		// Scan the result into variables
		err := rows.Scan(&status, &createdDate, &totalCount, &effortHours)
		if err != nil {
			return nil, fmt.Errorf("Error scanning ClickHouse result: %v", err)
		}

		result = append(result, bson.M{
			"status":      status,
			"createdDate": createdDate,
			"totalCount":  totalCount,
			"effortHours": effortHours,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating over rows: %v", err)
	}

	return result, nil
}
