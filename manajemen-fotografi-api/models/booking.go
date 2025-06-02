package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	BookingStatusPending   = "pending"
	BookingStatusConfirmed = "confirmed"
	BookingStatusDone      = "done"
)

type Booking struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ClientID       primitive.ObjectID `bson:"client_id" json:"client_id"`
	PhotographerID primitive.ObjectID `bson:"photographer_id" json:"photographer_id"`
	Date           time.Time          `bson:"date" json:"date"` // format ISO8601
	Location       string             `bson:"location" json:"location"`
	Status         string             `bson:"status" json:"status"` // gunakan konstanta
	Note           string             `bson:"note,omitempty" json:"note,omitempty"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
