package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"k8s-go-grpc-react/internal/logger"
	pb "k8s-go-grpc-react/proto"
)

const (
	httpPort = ":8081"
)

var (
	// Настройка логера с Graylog
	log = logger.SetupGraylogLogger("http-gateway")
)

func getGRPCServerAddr() string {
	if addr := os.Getenv("GRPC_SERVER_ADDR"); addr != "" {
		return addr
	}
	return "localhost:8080"
}

type Gateway struct {
	client pb.UserServiceClient
}

func NewGateway() (*Gateway, error) {
	grpcAddr := getGRPCServerAddr()

	log.WithFields(logrus.Fields{
		"component": "gateway-init",
		"grpc_addr": grpcAddr,
	}).Info("Подключение к gRPC серверу")

	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.WithError(err).Error("Не удалось подключиться к gRPC серверу")
		return nil, fmt.Errorf("не удалось подключиться к gRPC серверу: %v", err)
	}

	client := pb.NewUserServiceClient(conn)

	log.WithField("component", "gateway-init").Info("Успешно подключились к gRPC серверу")

	return &Gateway{client: client}, nil
}

func (g *Gateway) enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func (g *Gateway) handleOptions(w http.ResponseWriter, r *http.Request) {
	g.enableCORS(w)
	w.WriteHeader(http.StatusOK)
}

func (g *Gateway) extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func (g *Gateway) createAuthContext(ctx context.Context, token string) context.Context {
	if token == "" {
		return ctx
	}

	md := metadata.Pairs("authorization", "Bearer "+token)
	return metadata.NewOutgoingContext(ctx, md)
}

func (g *Gateway) getUser(w http.ResponseWriter, r *http.Request) {
	g.enableCORS(w)

	vars := mux.Vars(r)
	idStr := vars["id"]

	log.WithFields(logrus.Fields{
		"component": "get-user",
		"user_id":   idStr,
		"method":    r.Method,
		"path":      r.URL.Path,
	}).Info("Запрос получения пользователя")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithError(err).WithField("user_id", idStr).Error("Неверный ID пользователя")
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Добавляем токен в контекст
	token := g.extractToken(r)
	ctx = g.createAuthContext(ctx, token)

	resp, err := g.client.GetUser(ctx, &pb.GetUserRequest{Id: int32(id)})
	if err != nil {
		log.WithError(err).WithField("user_id", id).Error("Ошибка получения пользователя")
		http.Error(w, fmt.Sprintf("Ошибка получения пользователя: %v", err), http.StatusInternalServerError)
		return
	}

	log.WithFields(logrus.Fields{
		"component": "get-user",
		"user_id":   id,
		"user_name": resp.User.Name,
	}).Info("Пользователь успешно получен")

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.WithError(err).Error("Ошибка кодирования ответа")
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}

func (g *Gateway) createUser(w http.ResponseWriter, r *http.Request) {
	g.enableCORS(w)

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithError(err).Error("Неверный JSON в запросе создания пользователя")
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	log.WithFields(logrus.Fields{
		"component":  "create-user",
		"user_name":  req.Name,
		"user_email": req.Email,
	}).Info("Запрос создания пользователя")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Добавляем токен в контекст
	token := g.extractToken(r)
	ctx = g.createAuthContext(ctx, token)

	resp, err := g.client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"user_name":  req.Name,
			"user_email": req.Email,
		}).Error("Ошибка создания пользователя")
		http.Error(w, fmt.Sprintf("Ошибка создания пользователя: %v", err), http.StatusInternalServerError)
		return
	}

	log.WithFields(logrus.Fields{
		"component": "create-user",
		"user_id":   resp.User.Id,
		"user_name": resp.User.Name,
	}).Info("Пользователь успешно создан")

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.WithError(err).Error("Ошибка кодирования ответа")
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}

func (g *Gateway) listUsers(w http.ResponseWriter, r *http.Request) {
	g.enableCORS(w)

	log.WithFields(logrus.Fields{
		"component": "list-users",
		"method":    r.Method,
		"path":      r.URL.Path,
	}).Info("Запрос списка пользователей")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Добавляем токен в контекст
	token := g.extractToken(r)
	ctx = g.createAuthContext(ctx, token)

	resp, err := g.client.ListUsers(ctx, &pb.Empty{})
	if err != nil {
		log.WithError(err).Error("Ошибка получения списка пользователей")
		http.Error(w, fmt.Sprintf("Ошибка получения списка пользователей: %v", err), http.StatusInternalServerError)
		return
	}

	log.WithFields(logrus.Fields{
		"component":   "list-users",
		"users_count": len(resp.Users),
	}).Info("Список пользователей успешно получен")

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.WithError(err).Error("Ошибка кодирования ответа")
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}

func (g *Gateway) register(w http.ResponseWriter, r *http.Request) {
	g.enableCORS(w)

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithError(err).Error("Неверный JSON в запросе регистрации")
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	log.WithFields(logrus.Fields{
		"component":  "register",
		"user_name":  req.Name,
		"user_email": req.Email,
	}).Info("Запрос регистрации пользователя")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := g.client.Register(ctx, &pb.RegisterRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"user_name":  req.Name,
			"user_email": req.Email,
		}).Error("Ошибка регистрации пользователя")
		http.Error(w, fmt.Sprintf("Ошибка регистрации: %v", err), http.StatusBadRequest)
		return
	}

	log.WithFields(logrus.Fields{
		"component": "register",
		"user_id":   resp.User.Id,
		"user_name": resp.User.Name,
	}).Info("Пользователь успешно зарегистрирован")

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.WithError(err).Error("Ошибка кодирования ответа")
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}

func (g *Gateway) login(w http.ResponseWriter, r *http.Request) {
	g.enableCORS(w)

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.WithError(err).Error("Неверный JSON в запросе входа")
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	log.WithFields(logrus.Fields{
		"component":  "login",
		"user_email": req.Email,
	}).Info("Запрос входа пользователя")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := g.client.Login(ctx, &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"user_email": req.Email,
		}).Error("Ошибка входа пользователя")
		http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
		return
	}

	log.WithFields(logrus.Fields{
		"component": "login",
		"user_id":   resp.User.Id,
		"user_name": resp.User.Name,
	}).Info("Пользователь успешно вошел в систему")

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.WithError(err).Error("Ошибка кодирования ответа")
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}

func main() {
	log.WithField("component", "startup").Info("Запуск HTTP Gateway...")

	gateway, err := NewGateway()
	if err != nil {
		log.WithError(err).Fatal("Не удалось создать gateway")
	}

	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// v1 API routes
	v1 := api.PathPrefix("/v1").Subrouter()

	// Auth routes
	v1.HandleFunc("/auth/register", gateway.register).Methods("POST")
	v1.HandleFunc("/auth/login", gateway.login).Methods("POST")

	// User routes
	v1.HandleFunc("/users/{id:[0-9]+}", gateway.getUser).Methods("GET")
	v1.HandleFunc("/users", gateway.createUser).Methods("POST")
	v1.HandleFunc("/users", gateway.listUsers).Methods("GET")

	// Legacy API routes (for backward compatibility)
	api.HandleFunc("/users/{id:[0-9]+}", gateway.getUser).Methods("GET")
	api.HandleFunc("/users", gateway.createUser).Methods("POST")
	api.HandleFunc("/users", gateway.listUsers).Methods("GET")

	// CORS preflight
	api.PathPrefix("/").HandlerFunc(gateway.handleOptions).Methods("OPTIONS")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(logrus.Fields{
			"component": "health-check",
			"method":    r.Method,
			"path":      r.URL.Path,
		}).Info("Health check запрос")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.WithError(err).Error("Ошибка записи ответа health check")
		}
	})

	log.WithFields(logrus.Fields{
		"component": "http-server",
		"port":      httpPort,
	}).Info("HTTP Gateway запущен и готов к приему запросов")

	log.WithError(http.ListenAndServe(httpPort, r)).Fatal("HTTP сервер остановлен")
}
