package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Gallery struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PhotographerID primitive.ObjectID `bson:"photographer_id" json:"photographer_id"`
	Title          string             `bson:"title" json:"title"`
	ImageURL       string             `bson:"image_url" json:"image_url"` // bisa juga []string jika banyak gambar
	Description    string             `bson:"description,omitempty" json:"description,omitempty"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
