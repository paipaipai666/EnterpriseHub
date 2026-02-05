package main

import (
	"net"

	"github.com/paipaipai666/EnterpriseHub/payment-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/payment-service/internal/handler"
	"github.com/paipaipai666/EnterpriseHub/payment-service/internal/pb"
	"github.com/paipaipai666/EnterpriseHub/payment-service/internal/repository"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func init() {
	initializers.LoadEnv()
	initializers.InitLogger("payment_service")
	initializers.ConnectToRedis()
	initializers.ConnectToDatabase()
}

func main() {
	startGrpcServer()
}

func startGrpcServer() {
	paymentRepo := repository.NewPaymentRepository()
	paymentGrpcHandler := handler.NewPaymentGrpcHandler(paymentRepo)

	grpcServer := grpc.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, paymentGrpcHandler)

	address := "0.0.0.0:11001"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		initializers.Log.Fatal("cannot tcp start server:", zap.Error(err))
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		initializers.Log.Fatal("cannot gRPC start server:", zap.Error(err))
	}
}
