package repository

import (
	"context"
	"log/slog"

	usecase "github.com/VanLavr/auth/internal/auth/service"
	"github.com/VanLavr/auth/internal/models"
	"github.com/VanLavr/auth/internal/pkg/config"
	e "github.com/VanLavr/auth/internal/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type authRepository struct {
	conn       string
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

func New(cfg *config.Config) usecase.Repository {
	return &authRepository{conn: cfg.Mongo}
}

// Connet() connects to mongo and selects the database and the collection.
func (a *authRepository) Connect(ctx context.Context, cfg *config.Config) error {
	// Create client options.
	clientOptions := options.Client().ApplyURI(a.conn)
	// Connect.
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	if err = client.Ping(ctx, nil); err != nil {
		slog.Error(err.Error())
		return err
	}
	slog.Info("pinged")

	// Select database and collection.
	a.client = client
	a.database = client.Database(cfg.DBName)
	a.collection = a.database.Collection(cfg.CollectionName)

	return nil
}

func (a *authRepository) CloseConnetion(ctx context.Context) error {
	if err := a.client.Disconnect(ctx); err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

// Create a filter.
// Find a token via tokenstring and guid.
// Bind it to an object and check if it's fields empty or not.
func (a *authRepository) GetToken(ctx context.Context, provided models.RefreshToken) (*models.RefreshToken, error) {
	// Create a filter.
	filter := bson.M{
		"guid": provided.GUID,
	}

	// Find a token via tokenstring and guid.
	cursor, err := a.collection.Find(ctx, filter)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	// Bind it to an object and check if it's fields empty or not.
	var result models.RefreshToken
	for cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			slog.Error(err.Error())
			return nil, err
		}
	}

	if result.GUID == "" || result.TokenString == "" {
		return nil, e.ErrTokenNotFound
	}

	return &result, nil
}

// Store generated refresh token.
func (a *authRepository) StoreToken(ctx context.Context, token models.RefreshToken) error {
	_, err := a.collection.InsertOne(ctx, token)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

// Create a filter.
// Create an updated document.
// Update a document that matches the filter.
func (a *authRepository) UpdateToken(ctx context.Context, provided models.RefreshToken) error {
	// Create a filter.
	filter := bson.M{
		"guid": provided.GUID,
	}

	// Create an updated document.
	update := bson.M{
		"$set": bson.M{
			"tokenstring": provided.TokenString,
			"guid":        provided.GUID,
		},
	}

	// Update a document that matches the filter.
	result, err := a.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	if result.MatchedCount != 1 {
		slog.Error("no matches")
		return e.ErrInternal
	}

	return nil
}
