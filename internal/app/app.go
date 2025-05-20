package app

import (
	"context"
	"fmt"
	"github.com/stawwkom/auth_service/internal/interceptor"
	desc "github.com/stawwkom/auth_service/pkg/auth_v1"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"

	runtimes "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	api "github.com/stawwkom/auth_service/internal/api/auth"
	"github.com/stawwkom/auth_service/internal/config"
	"github.com/stawwkom/auth_service/internal/repository"
	repo "github.com/stawwkom/auth_service/internal/repository/auth"
	serv "github.com/stawwkom/auth_service/internal/service/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	log.Printf("gRPC сервер запущен на %s (уровень логирования: %s)",
		a.serviceProvider.GRPCAddr(), config.Cfg.Log)

	listener, err := net.Listen("tcp", a.serviceProvider.GRPCAddr())
	if err != nil {
		return err
	}

	// Запускаем gRPC сервер в отдельной горутине
	go func() {
		if err := a.grpcServer.Serve(listener); err != nil {
			log.Printf("gRPC сервер остановлен: %v", err)
		}
	}()

	// Запускаем HTTP сервер в отдельной горутине
	go func() {
		log.Printf("HTTP сервер запущен на %s", a.serviceProvider.HTTPAddr())
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP сервер остановлен: %v", err)
		}
	}()

	// Ожидаем завершения контекста
	<-ctx.Done()

	// Завершаем gRPC сервер
	log.Println("⏹ Останавливаем gRPC сервер...")
	a.grpcServer.GracefulStop()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	steps := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
	}

	for _, step := range steps {
		if err := step(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	return config.Load()
}

func (a *App) initServiceProvider(ctx context.Context) error {
	db, err := repository.NewPostgresDB(ctx, config.Cfg)
	if err != nil {
		return err
	}

	authRepo := repo.NewRepository(db)
	authService := serv.NewAuthService(authRepo)

	a.serviceProvider = &serviceProvider{
		authService: authService,
		dbClose: func() error {
			db.Close()
			return nil
		},
	}
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtimes.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := desc.RegisterUserAPIHandlerFromEndpoint(
		ctx, mux, a.serviceProvider.GRPCAddr(), opts,
	); err != nil {
		return err
	}

	a.httpServer = &http.Server{
		Addr:    a.serviceProvider.HTTPAddr(),
		Handler: mux,
	}

	return nil
}

func (a *App) CloseHTTP() {
	if a.serviceProvider.dbClose != nil {
		if err := a.serviceProvider.dbClose(); err != nil {
			log.Printf("Ошибка при закрытии БД: %v", err)
		}
	}
	if a.grpcServer != nil {
		log.Println("⏹ Остановка gRPC сервера...")
		a.grpcServer.GracefulStop()
	}
	if a.httpServer != nil {
		log.Println("⏹ Остановка HTTP сервера...")
		if err := a.httpServer.Shutdown(context.Background()); err != nil {
			log.Printf("Ошибка при остановке HTTP сервера: %v", err)
		}
	}
}

//func (a *App) initHTTPServer(ctx context.Context) error {
//	mux := runtime.NewServeMux()
//
//	opts := []grpc.DialOption{
//		grpc.WithTransportCredentials(insecure.NewCredentials()),
//	}
//
//	err := desc.RegisterUserAPIHandlerFromEndpoint(ctx, mux, a.serviceProvider.GRPCAddr(), opts)
//	if err != nil {
//		return err
//	}
//
//	a.httpServer = &http.Server{
//		Addr: a.serviceProvider.HTTPAddr(),
//		Handler: mux,
//	}
//}

func (a *App) initGRPCServer(_ context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor),
	)
	fmt.Println("Validate running")
	reflection.Register(a.grpcServer)

	desc.RegisterUserAPIServer(a.grpcServer, api.NewServer(a.serviceProvider.authService))

	return nil
}

func (a *App) Close() {
	if a.serviceProvider.dbClose != nil {
		if err := a.serviceProvider.dbClose(); err != nil {
			log.Printf("Ошибка при закрытии БД: %v", err)
		}
	}
	if a.grpcServer != nil {
		log.Println("⏹ Остановка gRPC сервера...")
		a.grpcServer.GracefulStop()
	}
	if a.httpServer != nil {
		log.Println("⏹ Остановка HTTP сервера...")
		if err := a.httpServer.Shutdown(context.Background()); err != nil {
			log.Printf("Ошибка при остановке HTTP сервера: %v", err)
		}
	}
}
