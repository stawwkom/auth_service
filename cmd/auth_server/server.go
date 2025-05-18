package main

import (
	"context"
	"fmt"
	api "github.com/stawwkom/auth_service/internal/api/auth"
	"github.com/stawwkom/auth_service/internal/config"
	"github.com/stawwkom/auth_service/internal/repository"
	repo "github.com/stawwkom/auth_service/internal/repository/auth"
	"github.com/stawwkom/auth_service/internal/service"
	serv "github.com/stawwkom/auth_service/internal/service/auth"
	"log"
	"net"

	pb "github.com/stawwkom/auth_service/pkg/auth_v1"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserAPIServer
	authService service.AuthService
}

func main() {
	// Создаем контекст с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Загружаем конфигурацию (local.yaml + .env + переменные окружения)
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Получаем конфигурацию из глобальной переменной
	cfg := config.Cfg

	// Инициализируем подключение к базе данных
	db, err := repository.NewPostgresDB(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Инициализируем репозиторий
	authRepo := repo.NewRepository(db)

	// Инициализируем сервис
	authService := serv.NewAuthService(authRepo)

	// Формируем адрес для запуска gRPC сервера: "host:port"
	address := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Ошибка запуска: %v", err)
	}

	// Создаем наш gRPC server
	s := grpc.NewServer()

	// Регистрируем gRPC
	pb.RegisterUserAPIServer(s, api.NewServer(authService))

	log.Printf("gRPC сервер запущен на %v (уровень логирования: %v)", address, cfg.Log)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
