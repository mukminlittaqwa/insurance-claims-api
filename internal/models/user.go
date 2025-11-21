package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Username string `bson:"username" json:"username"`
    Password string `bson:"password" json:"-"`
    Role     string `bson:"role" json:"role"` // user, verifier, approver
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    Token string `json:"token"`
    User  struct {
        ID   primitive.ObjectID `json:"id"`
        Username string `json:"username"`
        Role string `json:"role"`
    } `json:"user"`
}