package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Entry struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Dish        *string            `json:"dish"`
	Size        *uint64            `json:"size"`
	Ingredients *string            `json:"ingredients"`
	Proteins    *uint64            `json:"proteins"`
	Carbs       *uint64            `json:"carbs"`
	Fat         *uint64            `json:"fat"`
	Calories    *uint64            `json:"calories"`
	CreatedAt   *time.Time         `bson:"createdAt"`
	UpdatedAt   *time.Time         `bson:"updatedAt"`
	// for soft delete
	SoftDeletedAt *time.Time `bson:"softDeletedAt"`
}
