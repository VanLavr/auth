package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/VanLavr/auth/internal/auth/delivery"
	"github.com/VanLavr/auth/internal/auth/usecase"
	"github.com/VanLavr/auth/internal/pkg/config"
)

func main() {
	usecase := usecase.New()
	srv := delivery.New(&config.Config{Addr: ":3000", Secret: "afdjsalf"}, usecase)
	srv.BindRoutes()

	go func() {
		log.Println("listening on :3000")
		if err := srv.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	context, close := signal.NotifyContext(context.Background(), os.Interrupt)
	defer close()

	<-context.Done()

	if err := srv.ShutDown(context); err != nil {
		log.Fatal(err)
	}
}
