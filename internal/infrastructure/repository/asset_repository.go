package repository

import (
	"context"

	"github.com/mikiasgoitom/Secure-Asset/internal/contract"
	"github.com/mikiasgoitom/Secure-Asset/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type assetRepository struct {
	collection *mongo.Collection
}

func NewAssetRepository(db *mongo.Database, collectionName string) contract.IAssetRepository {
	return &assetRepository{
		collection: db.Collection(collectionName),
	}
}

func (r *assetRepository) Create(ctx context.Context, asset *entity.Asset) (*entity.Asset, error) {
	_, err := r.collection.InsertOne(ctx, asset)
	if err != nil {
		return nil, err
	}
	return asset, nil
}
func (r *assetRepository) GetByID(ctx context.Context, id string) (*entity.Asset, error) {
	var asset entity.Asset
	filter := bson.M{"id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&asset)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}
func (r *assetRepository) Update(ctx context.Context, asset *entity.Asset) (*entity.Asset, error) {
	filter := bson.M{"id": asset.ID}
	update := bson.M{"$set": asset}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return asset, nil
}