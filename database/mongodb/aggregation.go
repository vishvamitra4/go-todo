package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AggregationWithOptions(collection *mongo.Collection, matchFields, groupByFields, projectFields, sortOptions bson.D, limitOptions int) ([]interface{}, error) {
	pipeline := mongo.Pipeline{}

	// matching stage...
	if len(matchFields) > 0 {
		matchStage := bson.D{{Key: "$match", Value: matchFields}}
		pipeline = append(pipeline, matchStage)
	}

	// group stage....
	if len(groupByFields) > 0 {
		groupByFields = append(groupByFields, bson.E{Key: "totalEffortHours", Value: bson.D{{Key: "$sum", Value: "$effort_hours"}}})
		groupStage := bson.D{{Key: "$group", Value: groupByFields}}
		pipeline = append(pipeline, groupStage)
	} else {
		groupStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "totalEffortHours", Value: bson.D{{Key: "$sum", Value: "$effort_hours"}}},
		}}}
		pipeline = append(pipeline, groupStage)
	}

	if len(projectFields) > 0 {
		projectStage := bson.D{{Key: "$project", Value: projectFields}}
		pipeline = append(pipeline, projectStage)
	}

	if len(sortOptions) > 0 {
		sortStage := bson.D{{Key: "$sort", Value: sortOptions}}
		pipeline = append(pipeline, sortStage)
	}

	if limitOptions > 0 {
		limitStage := bson.D{{Key: "$limit", Value: limitOptions}}
		pipeline = append(pipeline, limitStage)
	}

	cur, err := collection.Aggregate(context.Background(), pipeline, options.Aggregate().SetMaxTime(30*time.Second).SetAllowDiskUse(true))
	var result []interface{}
	if err != nil {
		return result, err
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var doc interface{}
		err := cur.Decode(&doc)
		if err != nil {
			return result, err
		}
		result = append(result, doc)
	}

	return result, nil
}
