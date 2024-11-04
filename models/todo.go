package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Todo represents a todo item
type Todo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`         // MongoDB ObjectID
	Title       string             `bson:"title"`                 // Title of the task
	Desc        string             `bson:"description,omitempty"` // Description of the task
	Status      string             `bson:"status,"`               // Status of the task (e.g., "Pending", "Completed")
	CreatedAt   time.Time          `bson:"created"`               // Date and time of creation
	EffortHours int                `bson:"effort_hours"`          // Estimated effort hours required to complete the task
}
