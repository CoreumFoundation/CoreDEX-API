package main

import (
	"net"
	"os"

	"github.com/CoreumFoundation/CoreDEX-API/apps/store/ports/grpc"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

var (
	grpcPort string
	lis      net.Listener
)

// Initialize any functionality which would lead to application termination if not set correctly
// by using environment variables
func init() {
	parseConfig()
	var err error
	lis, err = net.Listen("tcp", grpcPort)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
}

// Start application
func main() {
	grpcServer := grpc.NewGrpcServer()
	logger.Infof("listening at %v", lis.Addr())
	if err := grpcServer.Server.Serve(lis); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}

// parses the provided config and checks if all required values are set
func parseConfig() {
	grpcPort = os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		logger.Fatalf("GRPC_PORT env is not set (format: :port, usually :50051)")
	}
}
