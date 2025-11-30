package repository

import (
	"context"

	"github.com/mikiasgoitom/Secure-Asset/internal/contract"
	"github.com/mikiasgoitom/Secure-Asset/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database, collectionName string) contract.IUserRepository {
	return &userRepository{
		collection: db.Collection(collectionName),
	}
}


// Create inserts a new user document into the 'users' collection.
func (r *userRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
    _, err := r.collection.InsertOne(ctx, user)
    if err != nil {
        return nil, err
    }
    return user, nil
}

// FindByEmail retrieves a user by their email address.
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
    var user entity.User
    err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, nil // Return nil, nil if no user is found
        }
        return nil, err
    }
    return &user, nil
}

// FindByUsername retrieves a user by their username.
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
    var user entity.User
    err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, nil // Return nil, nil if no user is found
        }
        return nil, err
    }
    return &user, nil
}

// FindByID retrieves a user by their unique ID.
func (r *userRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
    var user entity.User
    err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

// Update modifies an existing user's details in the database.
func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
    _, err := r.collection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})
    return err
}

