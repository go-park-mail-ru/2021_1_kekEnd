package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/proto"
	sessionsGrpc "github.com/go-park-mail-ru/2021_1_kekEnd/internal/services/sessions/delivery/grpc"
	sessionsRepo "github.com/go-park-mail-ru/2021_1_kekEnd/internal/services/sessions/repository"
	sessionsUC "github.com/go-park-mail-ru/2021_1_kekEnd/internal/services/sessions/usecase"
	constants "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost:%s", constants.RedisPort),
		Password: "",
		DB:       0,
	})

	p, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Failed to create redis client", p, err)
	}

	repo := sessionsRepo.NewRedisRepository(rdb)
	usecase := sessionsUC.NewUseCase(repo)
	handler := sessionsGrpc.NewAuthHandlerServer(usecase)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", constants.AuthPort))

	if err != nil {
		log.Fatalln("Can't listen session microservice port", err)
	}
	defer lis.Close()

	server := grpc.NewServer()
	proto.RegisterAuthHandlerServer(server, handler)

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
