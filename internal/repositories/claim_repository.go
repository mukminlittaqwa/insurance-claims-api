package repositories

import (
    "context"
    "insurance-claims-api/internal/models"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type ClaimRepository interface {
    Create(claim *models.Claim) error
    FindByID(id primitive.ObjectID) (*models.Claim, error)
    FindByUserID(userID primitive.ObjectID, page, limit int) ([]models.Claim, int64, error)
    FindAll(filter bson.M, page, limit int) ([]models.Claim, int64, error)
    Update(claim *models.Claim) error
    Delete(id primitive.ObjectID) error
    AddHistory(id primitive.ObjectID, history models.ClaimHistory) error
    UpdateWithPush(id primitive.ObjectID, update bson.M) error
}

type claimRepository struct {
    collection *mongo.Collection
}

func NewClaimRepository(client *mongo.Client) ClaimRepository {
    collection := client.Database("insurance").Collection("claims")
    return &claimRepository{collection}
}

func (r *claimRepository) Create(claim *models.Claim) error {
    claim.CreatedAt = time.Now()
    claim.UpdatedAt = time.Now()
    claim.Status = models.Draft
    claim.History = []models.ClaimHistory{
        {Status: models.Draft, ChangedBy: claim.UserID, ChangedAt: time.Now()},
    }
    _, err := r.collection.InsertOne(context.TODO(), claim)
    return err
}

func (r *claimRepository) FindByID(id primitive.ObjectID) (*models.Claim, error) {
    var claim models.Claim
    err := r.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&claim)
    if err != nil {
        return nil, err
    }
    return &claim, nil
}

func (r *claimRepository) FindByUserID(userID primitive.ObjectID, page, limit int) ([]models.Claim, int64, error) {
    skip := (page - 1) * limit
    opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.M{"created_at": -1})
    cursor, err := r.collection.Find(context.TODO(), bson.M{"user_id": userID}, opts)
    if err != nil {
        return nil, 0, err
    }
    defer cursor.Close(context.TODO())

    var claims []models.Claim
    if err = cursor.All(context.TODO(), &claims); err != nil {
        return nil, 0, err
    }
    total, _ := r.collection.CountDocuments(context.TODO(), bson.M{"user_id": userID})
    return claims, total, nil
}

func (r *claimRepository) FindAll(filter bson.M, page, limit int) ([]models.Claim, int64, error) {
    skip := (page - 1) * limit
    opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.M{"created_at": -1})
    cursor, err := r.collection.Find(context.TODO(), filter, opts)
    if err != nil {
        return nil, 0, err
    }
    defer cursor.Close(context.TODO())

    var claims []models.Claim
    if err = cursor.All(context.TODO(), &claims); err != nil {
        return nil, 0, err
    }
    total, _ := r.collection.CountDocuments(context.TODO(), filter)
    return claims, total, nil
}

func (r *claimRepository) Update(claim *models.Claim) error {
    claim.UpdatedAt = time.Now()
    _, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": claim.ID},
        bson.M{"$set": claim})
    return err
}

func (r *claimRepository) Delete(id primitive.ObjectID) error {
    _, err := r.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
    return err
}

func (r *claimRepository) AddHistory(id primitive.ObjectID, history models.ClaimHistory) error {
    _, err := r.collection.UpdateOne(context.TODO(),
        bson.M{"_id": id},
        bson.M{
            "$push": bson.M{"history": history},
            "$set":  bson.M{"updated_at": time.Now()},
        },
    )
    return err
}

func (r *claimRepository) UpdateWithPush(id primitive.ObjectID, update bson.M) error {
    _, err := r.collection.UpdateOne(
        context.TODO(),
        bson.M{"_id": id},
        update,
    )
    return err
}