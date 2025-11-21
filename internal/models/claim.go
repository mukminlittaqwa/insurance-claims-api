package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type ClaimStatus string

const (
    Draft     ClaimStatus = "draft"
    Submitted ClaimStatus = "submitted"
    Reviewed  ClaimStatus = "reviewed"
    Approved  ClaimStatus = "approved"
    Rejected  ClaimStatus = "rejected"
)

type ClaimHistory struct {
    Status    ClaimStatus        `bson:"status" json:"status"`
    ChangedBy primitive.ObjectID `bson:"changed_by" json:"changed_by"`
    ChangedAt time.Time          `bson:"changed_at" json:"changed_at"`
    Note      string             `bson:"note,omitempty" json:"note,omitempty"`
}

type Claim struct {
    ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
    PolicyNumber string             `bson:"policy_number" json:"policy_number" binding:"required"`
    ClaimAmount  float64            `bson:"claim_amount" json:"claim_amount" binding:"required"`
    Description  string             `bson:"description" json:"description" binding:"required"`
    Documents    []string           `bson:"documents,omitempty" json:"documents,omitempty"`
    Status       ClaimStatus        `bson:"status" json:"status"`
    History      []ClaimHistory     `bson:"history" json:"history"`
    CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateClaimRequest struct {
    PolicyNumber string   `json:"policy_number" binding:"required"`
    ClaimAmount  float64  `json:"claim_amount" binding:"required"`
    Description  string   `json:"description" binding:"required"`
    Documents    []string `json:"documents,omitempty"`
}

type UpdateClaimRequest struct {
    PolicyNumber string   `json:"policy_number,omitempty"`
    ClaimAmount  float64  `json:"claim_amount,omitempty"`
    Description  string   `json:"description,omitempty"`
    Documents    []string `json:"documents,omitempty"`
}