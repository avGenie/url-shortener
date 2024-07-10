package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/avGenie/url-shortener/internal/app/config"
	"github.com/avGenie/url-shortener/internal/app/grpc/interceptor"
	handlers "github.com/avGenie/url-shortener/internal/app/handlers/delete"
	storage_api "github.com/avGenie/url-shortener/internal/app/storage/api/model"
	pb "github.com/avGenie/url-shortener/proto"
	"google.golang.org/grpc"
)

// ShortenerServer GRPC server
type ShortenerServer struct {
	pb.ShortenerServer

	storage storage_api.Storage
	config  config.Config

	server        *grpc.Server
	deleteHandler *handlers.DeleteHandler
}

// NewGRPCServer Creates new GRPC server
func NewGRPCServer(config config.Config, storage storage_api.Storage) *ShortenerServer {
	return &ShortenerServer{
		storage:       storage,
		config:        config,
		server:        grpc.NewServer(grpc.UnaryInterceptor(interceptor.AuthInterceptor)),
		deleteHandler: handlers.NewDeleteHandler(storage),
	}
}

// Start Starts GRPC server
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

// Stop Stops GRPC server
func (s *ShortenerServer) Stop() {
	s.server.GracefulStop()
}
