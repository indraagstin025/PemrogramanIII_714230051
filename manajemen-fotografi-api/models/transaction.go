package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	TransactionStatusPaid   = "paid"
	TransactionStatusUnpaid = "unpaid"
)

type Transaction struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookingID  primitive.ObjectID `bson:"booking_id" json:"booking_id"`
	Method     string             `bson:"method" json:"method"`         // contoh: "transfer", "ewallet"
	Total      float64            `bson:"total" json:"total"`           // misal: 500000
	Status     string             `bson:"status" json:"status"`         // contoh: "paid", "unpaid"
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"` // opsional
}
