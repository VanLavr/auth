package repository

import (
	"context"
	"log"

	"github.com/VanLavr/auth/internal/auth/usecase"
	"github.com/VanLavr/auth/internal/models"
	"github.com/VanLavr/auth/internal/pkg/config"
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

func (a *authRepository) GetToken(ctx context.Context, provided models.RefreshToken) (*models.RefreshToken, error) {
	cursor, err := a.collection.Find(ctx, bson.D{
		{"Token_String", provided.TokenString},
		{"GUID", provided.GUID},
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var result models.RefreshToken
	for cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (a *authRepository) StoreToken(ctx context.Context, token models.RefreshToken) error {
	_, err := a.collection.InsertOne(ctx, token)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (a *authRepository) MarkToken(models.RefreshToken) error {
	panic("not implemented")
}
