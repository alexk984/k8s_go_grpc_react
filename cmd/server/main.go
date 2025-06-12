package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"k8s-go-grpc-react/internal/auth"
	"k8s-go-grpc-react/internal/config"
	"k8s-go-grpc-react/internal/database"
	"k8s-go-grpc-react/internal/repository"
	"k8s-go-grpc-react/internal/service"
	pb "k8s-go-grpc-react/proto"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Подключаемся к базе данных используя отдельные переменные окружения
	dbConfig := database.GetConfigFromEnv()
	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	// Выполняем миграции
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Ошибка выполнения миграций: %v", err)
	}

	// Создаем репозиторий
	userRepo := repository.NewUserRepository(db)

	// Создаем метрики Prometheus
	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Общее количество gRPC запросов",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "Время выполнения gRPC запросов",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	usersCount := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "users_total",
			Help: "Общее количество пользователей",
		},
	)

	// Регистрируем метрики
	prometheus.MustRegister(requestsTotal, requestDuration, usersCount)

	// Создаем сервис с метриками
	userService := service.NewUserServiceWithMetrics(userRepo, requestsTotal, requestDuration, usersCount)

	// Создаем middleware для аутентификации
	authMiddleware := auth.NewAuthMiddleware()

	// Создаем gRPC сервер с middleware
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authMiddleware.UnaryInterceptor),
	)

	// Регистрируем сервис
	pb.RegisterUserServiceServer(grpcServer, userService)

	// Включаем рефлексию для отладки
	reflection.Register(grpcServer)

	// Запускаем gRPC сервер
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
		if err != nil {
			log.Fatalf("Ошибка создания TCP слушателя: %v", err)
		}

		log.Printf("gRPC сервер запущен на порту %s", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Ошибка запуска gRPC сервера: %v", err)
		}
	}()

	// Создаем HTTP сервер для gRPC-Gateway
	go func() {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		// Создаем gRPC-Gateway mux
		mux := runtime.NewServeMux()

		// Подключаемся к gRPC серверу
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		err := pb.RegisterUserServiceHandlerFromEndpoint(
			ctx,
			mux,
			fmt.Sprintf("localhost:%s", cfg.GRPCPort),
			opts,
		)
		if err != nil {
			log.Fatalf("Ошибка регистрации gRPC-Gateway: %v", err)
		}

		// Создаем HTTP роутер
		httpMux := http.NewServeMux()

		// Добавляем CORS middleware
		corsHandler := func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

				if r.Method == "OPTIONS" {
					w.WriteHeader(http.StatusOK)
					return
				}

				h.ServeHTTP(w, r)
			})
		}

		// Регистрируем маршруты
		httpMux.Handle("/api/", http.StripPrefix("/api", corsHandler(mux)))
		httpMux.Handle("/metrics", promhttp.Handler())
		httpMux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		}))

		log.Printf("HTTP сервер запущен на порту %s", cfg.HTTPPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.HTTPPort), httpMux); err != nil {
			log.Fatalf("Ошибка запуска HTTP сервера: %v", err)
		}
	}()

	// Ожидаем сигнал завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Получен сигнал завершения, останавливаем сервер...")

	// Graceful shutdown
	grpcServer.GracefulStop()

	log.Println("Сервер остановлен")
}
