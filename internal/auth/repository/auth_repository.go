package repository

import (
	"context"
	"log"

	"github.com/VanLavr/auth/internal/auth/usecase"
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

func New(conn string) usecase.Repository {
	return &authRepository{conn: conn}
}

func (a *authRepository) Connect(ctx context.Context, cfg *config.Config) error {
	clientOptions := options.Client().ApplyURI(a.conn)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return err
	}

	a.client = client
	a.database = client.Database(cfg.DBName)
	a.collection = a.database.Collection(cfg.CollectionName)

	return nil
}

func (a *authRepository) CloseConnetion(ctx context.Context) error {
	if err := a.client.Disconnect(ctx); err != nil {
		return err
	}
	return nil
}

// Find a token via tokenstring and guid.
// Bind it to an object and check if it's fields empty or not.
func (a *authRepository) GetToken(ctx context.Context, provided models.RefreshToken) (*models.RefreshToken, error) {
	// Find a token via tokenstring and guid.
	cursor, err := a.collection.Find(ctx, bson.D{
		{"Token_String", provided.TokenString},
		{"GUID", provided.GUID},
	})

	if err != nil {
		return nil, err
	}

	// Bind it to an object and check if it's fields empty or not.
	var result models.RefreshToken
	for cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
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
		log.Println(err)
		return err
	}

	return nil
}

func (a *authRepository) UpdateToken(models.RefreshToken) error {
	panic("not implemented")
}
