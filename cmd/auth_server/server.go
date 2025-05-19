package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/stawwkom/auth_service/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}
	defer a.Close() // ✅ Закрытие ресурсов после завершения

	// Обработка SIGINT / SIGTERM
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		log.Println("⛔️ Signal received, shutting down gracefully...")
		cancel()
	}()

	if err := a.Run(ctx); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

}
