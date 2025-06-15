package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"project-name/internal/config"
	"project-name/internal/models"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		collection: config.DB.Collection("users"),
	}
}

func (r *UserRepository) Create(user *models.User) (*models.User, error) {
	user.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(config.MongoCtx, user)
	return user, err
}

func (r *UserRepository) FindByID(id string) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.collection.FindOne(config.MongoCtx, bson.M{"_id": objectID}).Decode(&user)
	return &user, err
}