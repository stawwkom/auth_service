package main

import (
	"context"
	"google.golang.org/grpc/credentials"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pb "github.com/stawwkom/auth_service/pkg/auth_v1"
)

const (
	certPath = "../../certs/service.pem"
)

func main() {
	// Устанавливаем соединение с gRPC сервером

	creds, err := credentials.NewClientTLSFromFile(certPath, "")
	if err != nil {
		log.Fatalf("could not process the credentials: %v", err)
	}

	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("не удалось подключиться: %v", err)
	}
	defer conn.Close()

	// Создаём клиента UserService
	client := pb.NewUserAPIClient(conn)

	// Контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// === Вызов Create ===
	createRes, err := client.Create(ctx, &pb.CreateUserRequest{
		Name:            "Alice",
		Email:           "alice@example.com",
		Password:        "secret",
		PasswordConfirm: "secret",
		Role:            pb.Role_ADMIN,
	})
	if err != nil {
		log.Fatalf("Create вызвал ошибку: %v", err)
	}
	log.Printf("Создан пользователь с ID: %d", createRes.Id)

	// === Вызов Get ===
	getRes, err := client.Get(ctx, &pb.GetUserRequest{Id: createRes.Id})
	if err != nil {
		log.Fatalf("Get вызвал ошибку: %v", err)
	}
	log.Printf("Получен пользователь: %+v", getRes)

	// === Вызов Update ===
	_, err = client.Update(ctx, &pb.UpdateUserRequest{
		Id:    createRes.Id,
		Name:  wrapperspb.String("Alice Updated"),
		Email: wrapperspb.String("alice.updated@example.com"),
	})
	if err != nil {
		log.Fatalf("Update вызвал ошибку: %v", err)
	}
	log.Println("Пользователь обновлён")

	// === Вызов Delete ===
	_, err = client.Delete(ctx, &pb.DeleteUserRequest{Id: createRes.Id})
	if err != nil {
		log.Fatalf("Delete вызвал ошибку: %v", err)
	}
	log.Println("Пользователь удалён")
}
