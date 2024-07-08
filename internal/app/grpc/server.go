package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/avGenie/url-shortener/internal/app/grpc/interceptor"
	storage_api "github.com/avGenie/url-shortener/internal/app/storage/api/model"
	pb "github.com/avGenie/url-shortener/proto"
	"google.golang.org/grpc"
)

type ShortenerServer struct {
	pb.ShortenerServer

	storage storage_api.Storage
	config  config.Config

	server *grpc.Server
}

func NewGRPCServer(config config.Config, storage storage_api.Storage) *ShortenerServer {
	return &ShortenerServer{
		storage: storage,
		config:  config,
		server: grpc.NewServer(grpc.UnaryInterceptor(interceptor.AuthInterceptor)),
	}
}

func (s *ShortenerServer) Start() {
	listen, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}

	// регистрируем сервис
	pb.RegisterShortenerServer(s.server, s)

	fmt.Println("Сервер gRPC начал работу")
	
	if err := s.server.Serve(listen); err != nil {
		log.Fatal(err)
	}
}

func (s *ShortenerServer) Stop() {
	s.server.GracefulStop()
}
