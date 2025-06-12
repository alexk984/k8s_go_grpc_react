#!/bin/bash

# Создаем директорию для сгенерированных файлов
mkdir -p proto

# Генерируем Go код из protobuf с поддержкой gRPC-Gateway
/opt/homebrew/bin/protoc \
    --proto_path=. \
    --proto_path=third_party/googleapis \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
    proto/user.proto

echo "Protobuf файлы успешно сгенерированы!" 