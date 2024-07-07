package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdatePostRequest struct {
	UserID      primitive.ObjectID `json:"user_id,omitempty"`
    Title       *string `json:"title,omitempty"`
    Description *string `json:"description,omitempty"`
    Status      *string `json:"status,omitempty"`
}

type DeletePostRequest struct {
	UserID      primitive.ObjectID `json:"user_id,omitempty"`
}

type Post struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    UserID      primitive.ObjectID `bson:"user_id,omitempty" json:"user_id"`
    Title       string             `bson:"title" json:"title"`
    Description string             `bson:"description" json:"description"`
    Status      string             `bson:"status" json:"status"`
    CreatedAt   time.Time          `bson:"created" json:"created"`
    UpdatedAt   time.Time          `bson:"updated" json:"updated"`
}
