package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"

	"grpc-gateway-test/auth"
	"grpc-gateway-test/env"
	pb "grpc-gateway-test/model"
	"grpc-gateway-test/service"
)

func main() {
	ctx := context.Background()
	go startGrpcServer(ctx)
	startGrpcGatewayServer(ctx)
}

func startGrpcServer(ctx context.Context) {

	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_auth.UnaryServerInterceptor(auth.AuthInterceptor),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	pb.RegisterDatasourceRpcServer(server, &service.DatasourceRpcService{})
	pb.RegisterUserRpcServer(server, &service.UserRpcService{})

	log.Printf("Starting gRPC server on %s...", env.GrpcPort)
	listenPort, err := net.Listen("tcp", fmt.Sprintf(":%s", env.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	err = server.Serve(listenPort)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	log.Printf("Start gRPC server on %s", env.GrpcPort)
}

func startGrpcGatewayServer(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	log.Printf("Starting gRPC gateway server on %s...", env.HttpPort)
	log.Printf("Connectiong gRPC server with %s...", env.GrpcPort)
	err := pb.RegisterDatasourceRpcHandlerFromEndpoint(ctx, mux,
		fmt.Sprintf("localhost:%s", env.GrpcPort), opts)
	if err != nil {
		log.Fatalf("failed to setup grpc gateway: %v", err)
	}

	err = pb.RegisterUserRpcHandlerFromEndpoint(ctx, mux,
		fmt.Sprintf("localhost:%s", env.GrpcPort), opts)
	if err != nil {
		log.Fatalf("failed to setup grpc gateway: %v", err)
	}

	err = http.ListenAndServe(fmt.Sprintf(":%s", env.HttpPort), mux)
	if err != nil {
		log.Fatalf("failed to start grpc gateway: %v", err)
	}
	log.Printf("Start gRPC gateway server on %s", env.HttpPort)
}
