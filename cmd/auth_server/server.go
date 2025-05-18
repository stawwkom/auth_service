package main

import (
	"context"
	"fmt"
	"github.com/stawwkom/auth_service/internal/config"
	"github.com/stawwkom/auth_service/internal/service"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/VictorSidoruk/auth/pkg"
)

type server struct {
	pb.UnimplementedUserAPIServer
	authService service.AuthService
}

func (s *server) Create(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	fmt.Printf("CreateUser: %+v\n", req)
	return &pb.CreateUserResponse{Id: 1}, nil
}

func (s *server) Get(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	fmt.Printf("GetUser: %+v\n", req)
	return &pb.GetUserResponse{
		Id:        req.Id,
		Name:      "stawwkom",
		Email:     "stawwkom@gmail.com",
		Role:      pb.Role_USER,
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}, nil
}

func (s *server) Update(ctx context.Context, req *pb.UpdateUserRequest) (*emptypb.Empty, error) {
	fmt.Printf("UpdateUser: %+v\n", req)
	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	fmt.Printf("DeleteUser: %+v\n", req)
	return &emptypb.Empty{}, nil
}

func main() {
	// Загружаем конфигурацию (local.yaml + .env + переменные окружения)
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Получаем конфигурацию из глобальной переменной
	cfg := config.Cfg

	// Формируем адрес для запуска gRPC сервера: "host:port"
	address := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Ошибка запуска: %v", err)
	}
	// Создаем наш gRPC server
	s := grpc.NewServer()

	// Регистрируем gRPC

	pb.RegisterUserAPIServer(s, &server{})

	log.Printf("gRPC сервер запущен на %v (уровень логирования: %v)", address, cfg.Log)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
