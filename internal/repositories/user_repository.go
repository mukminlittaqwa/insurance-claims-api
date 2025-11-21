package repositories

import (
    "context"
    "insurance-claims-api/internal/models"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
    FindByUsername(username string) (*models.User, error)
    Create(user *models.User) error
    FindByID(id primitive.ObjectID) (*models.User, error)
}

type userRepository struct {
    collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client) UserRepository {
    collection := client.Database("insurance").Collection("users")
    return &userRepository{collection}
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
    var user models.User
    err := r.collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) Create(user *models.User) error {
    _, err := r.collection.InsertOne(context.TODO(), user)
    return err
}

func (r *userRepository) FindByID(id primitive.ObjectID) (*models.User, error) {
    var user models.User
    err := r.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}