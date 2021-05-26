package main

import (
	"fmt"
	"log"
	"net"

	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/proto"
	fileServerGrpc "github.com/go-park-mail-ru/2021_1_kekEnd/internal/services/fileserver/delivery/grpc"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"google.golang.org/grpc"
)

func main() {
	handler := fileServerGrpc.NewFileServerHandlerServer()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", _const.FileServerPort))

	if err != nil {
		log.Fatalln("Can't listen session microservice port", err)
	}
	defer lis.Close()

	server := grpc.NewServer()
	proto.RegisterFileServerHandlerServer(server, handler)

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
