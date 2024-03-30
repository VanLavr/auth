package usecase_test

import (
	"context"
	"log"
	"testing"
	"time"

	usecase "github.com/VanLavr/auth/internal/auth/service"
	"github.com/VanLavr/auth/internal/pkg/config"
	e "github.com/VanLavr/auth/internal/pkg/errors"
	"github.com/VanLavr/auth/mocks"
)

func TestGetNewTokenPair(t *testing.T) {
	testCase := struct {
		ctx           context.Context
		id            string
		ExpectedPair  map[string]any
		ExpectedError error
	}{
		ctx:           context.TODO(),
		id:            "afoweff",
		ExpectedPair:  nil,
		ExpectedError: e.ErrInvalidGUID,
	}

	usecase := usecase.New(&mocks.Repository{}, &config.Config{
		Secret:         "asdf",
		AccessExpTime:  time.Second * 120,
		RefreshExpTime: time.Second * 240,
	})

	tokens, err := usecase.GetNewTokenPair(testCase.ctx, testCase.id)

	if err != testCase.ExpectedError || tokens != nil {
		log.Fatal("aaa")
	}
}
