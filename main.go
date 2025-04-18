package main

import (
	"backend/main/config"

	"backend/main/controllers"
	"context"

	"fmt"

	"google.golang.org/grpc"
	"backend/main/pb"
	"net"
)

func main() {
	config.ConnectPostgreSQL()
	config.InitFirebase()
	config.InitLogger()
	config.InitGenAI()
	ctx := context.Background()
	_, err := config.FirebaseApp.Database(ctx)
	if err != nil {
		return
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println("Failed to listen:", err)
		return
	}

	grpcServer := grpc.NewServer()

	userFeedService := controllers.InitUserFeedController()
	pb.RegisterUserFeedServiceServer(grpcServer, userFeedService)

	fmt.Println("gRPC server running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		fmt.Println("failed to serve: %v", err)
	}
}