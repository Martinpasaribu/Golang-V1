package userRepository

import (
	"context"
	"errors"

	"github.com/Martinpasaribu/Golang-V1/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		collection: db.Collection("users"),
	}
}

func (r *userRepository) CreateUser(user *models.User) (*models.User, error) {
	ctx := context.TODO()
	
	// Cek apakah email sudah terdaftar
	existingUser, err := r.FindUserByEmail(user.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Set timestamps
	user.BeforeCreate()

	// Insert user baru
	res, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (r *userRepository) FindUserByEmail(email string) (*models.User, error) {
	ctx := context.TODO()
	
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User tidak ditemukan bukan error
		}
		return nil, err
	}
	
	return &user, nil
}