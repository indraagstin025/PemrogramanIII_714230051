package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// RoleType membatasi hanya dua role yang diizinkan
const (
    RoleClient      = "client"
    RolePhotographer = "photographer"
)

// User mewakili akun login pengguna dengan role tertentu.
type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name      string             `bson:"name" json:"name"`
    Email     string             `bson:"email" json:"email"`
    Password  string             `bson:"password" json:"-"` // disembunyikan dari response JSON
    Role      string             `bson:"role" json:"role"`  // hanya "client" atau "photographer"
    CreatedAt int64              `bson:"created_at" json:"created_at"`
    UpdatedAt int64              `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

// Client menyimpan detail tambahan untuk user dengan role "client".
type Client struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
    Name      string             `bson:"name" json:"name"`          // field baru
    Phone     string             `bson:"phone" json:"phone"`
    Address   string             `bson:"address,omitempty" json:"address,omitempty"`
    CreatedAt int64              `bson:"created_at" json:"created_at"`
    UpdatedAt int64              `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}


// Photographer menyimpan detail tambahan untuk user dengan role "photographer".
type Photographer struct {
    ID           primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    UserID       primitive.ObjectID   `bson:"user_id,omitempty" json:"user_id"`
    Phone        string               `bson:"phone" json:"phone"`
    Description  string               `bson:"description" json:"description"`
    Portfolio    []string             `bson:"portfolio" json:"portfolio"`
    Location     string               `bson:"location" json:"location"`
    ProfilePhoto string               `bson:"profile_photo" json:"profile_photo"` // URL path ke foto profil
    CreatedAt    int64                `bson:"created_at" json:"created_at"`
    UpdatedAt    int64                `bson:"updated_at" json:"updated_at"`
}

