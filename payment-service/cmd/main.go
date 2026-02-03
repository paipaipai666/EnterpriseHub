package main

import (
	"log"
	"net"

	"github.com/paipaipai666/EnterpriseHub/payment-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/payment-service/internal/handler"
	"github.com/paipaipai666/EnterpriseHub/payment-service/internal/pb"
	"github.com/paipaipai666/EnterpriseHub/payment-service/internal/repository"
	"google.golang.org/grpc"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDatabase()
}

var (
	paymentRepo        repository.PaymentRepository = repository.NewPaymentRepository()
	paymentGrpcHandler handler.PaymentGrpcHandler   = *handler.NewPaymentGrpcHandler(paymentRepo)
)

func main() {
	startGrpcServer()
}

func startGrpcServer() {
	grpcServer := grpc.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, &paymentGrpcHandler)

	address := "0.0.0.0:11001"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
