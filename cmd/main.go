package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/VanLavr/auth/internal/auth/delivery"
	"github.com/VanLavr/auth/internal/pkg/config"
)

type s struct{}

func (this *s) RefreshToken(string) (map[string]string, error)
func (this *s) CheckIfTokenIsUsed(string, string) bool

func main() {
	srv := delivery.New(&config.Config{Addr: ":3000", Secret: "afdjsalf"}, &s{})
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
