package models

import (
	"time"
)

// Todo represents a todo item
type Todo struct {
	ID          interface{} `json:"id" bson:"_id"`         // MongoDB ObjectID and clickhouse uuid
	Title       string      `bson:"title"`                 // Title of the task
	Desc        string      `bson:"description,omitempty"` // Description of the task
	Status      string      `bson:"status,"`               // Status of the task (e.g., "Pending", "Completed")
	CreatedAt   time.Time   `bson:"created"`               // Date and time of creation
	EffortHours int         `bson:"effort_hours"`          // Estimated effort hours required to complete the task
}
