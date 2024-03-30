package repository

import (
	"context"

	"github.com/VanLavr/auth/internal/auth/usecase"
	"github.com/VanLavr/auth/internal/models"
	"github.com/VanLavr/auth/internal/pkg/config"
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

func (a *authRepository) StoreToken(token models.RefreshToken)

func (a *authRepository) Connect(cfg *config.Config) error {
	clientOptions := options.Client().ApplyURI(a.conn)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return err
	}

	a.client = client
	a.database = client.Database(cfg.DBName)
	a.collection = a.database.Collection(cfg.CollectionName)

	return nil
}

func (a *authRepository) CloseConnetion() error {
	if err := a.client.Disconnect(context.Background()); err != nil {
		return err
	}
	return nil
}

func (a *authRepository) SetCollection()

func (a *authRepository) GetToken(string) *models.RefreshToken
