package config


import "go.mongodb.org/mongo-driver/mongo"

var (
    BookingCollection     *mongo.Collection
    TransactionCollection *mongo.Collection
)

func InitCollections(db *mongo.Database) {
    BookingCollection = db.Collection("bookings")
    TransactionCollection = db.Collection("transactions")
}

func GetBookingCollection() *mongo.Collection {
    return BookingCollection
}

func GetTransactionCollection() *mongo.Collection {
    return TransactionCollection
}
