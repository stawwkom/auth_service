package main

import (
	"context"
	"github.com/stawwkom/auth_service/internal/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"

	"github.com/stawwkom/auth_service/internal/app"
)

func main() {
	// 1. Инициализация логгера
	err := logger.InitLogger("logs/app.log", true) // true — для dev (читаемый), false — prod (json)
	if err != nil {
		panic(err)
	}

	logger.Info("🚀 Старт приложения")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a, err := app.NewApp(ctx)
	if err != nil {
		logger.Fatal("❌ Ошибка при инициализации приложения", zap.Error(err))
	}
	defer a.Close()

	// 4. Обработка завершения
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		logger.Warn("⛔️ Сигнал завершения получен")
		cancel()
	}()

	// 5. Запуск
	if err := a.Run(ctx); err != nil {
		logger.Fatal("❌ Ошибка при запуске сервера", zap.Error(err))
	}

}
