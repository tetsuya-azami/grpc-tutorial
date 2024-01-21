package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	hellopb "mygrpc/pkg/grpc/hello"
	orderpb "mygrpc/pkg/grpc/order"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func NewMyServer() *myServer {
	return &myServer{}
}

func main() {
	// 1. 8080番portのLisnterを作成
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		panic(err)
	}

	// 2. gRPCサーバーを作成
	s := grpc.NewServer()
	hellopb.RegisterGreetingServiceServer(s, NewMyServer())
	orderpb.RegisterOrderServiceServer(s, NewMyServer())

	reflection.Register(s)

	// 3. 作成したgRPCサーバーを、8080番ポートで稼働させる
	go func() {
		log.Printf("start gRPC server port:%v", port)
		s.Serve(listener)
	}()

	fmt.Println("Server listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)

	// 4.Ctrl+Cが入力されたらGraceful shutdownされるようにする
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}

type myServer struct {
	hellopb.UnimplementedGreetingServiceServer
	orderpb.UnimplementedOrderServiceServer
}

// Unary RPC
func (s *myServer) Hello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	// リクエストからnameフィールドを取り出して
	// "Hello, [名前]!"というレスポンスを返す
	return &hellopb.HelloResponse{
		Message: fmt.Sprintf("Hello %s", req.GetName()),
	}, nil
}

// Server Stream RPC
func (s *myServer) HelloServerStream(req *hellopb.HelloRequest, stream hellopb.GreetingService_HelloServerStreamServer) error {
	resCount := 5

	// 5回レスポンスを返す
	for i := 0; i < resCount; i++ {
		// streamのSendメソッドでレスポンスを送信
		if err := stream.Send(&hellopb.HelloResponse{
			Message: fmt.Sprintf("[%d]Hello, %s!", i, req.GetName()),
		}); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}

	return nil
}

// Client Stream RPC
func (s *myServer) HelloClientStream(server hellopb.GreetingService_HelloClientStreamServer) error {
	nameList := make([]string, 0)
	for {
		req, err := server.Recv()
		if errors.Is(err, io.EOF) {
			message := fmt.Sprintf("Hello, %v!", nameList)
			return server.SendAndClose(&hellopb.HelloResponse{
				Message: message,
			})
		}
		if err != nil {
			return err
		}
		nameList = append(nameList, req.GetName())
	}
}

func (s *myServer) HelloBiStream(stream hellopb.GreetingService_HelloBiStreamServer) error {
	for {
		req, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			fmt.Println("all requests have been received")
			return nil
		}

		if err != nil {
			return err
		}

		if err := stream.Send(&hellopb.HelloResponse{
			Message: fmt.Sprintf("Hello, %s!", req.GetName()),
		}); err != nil {
			return err
		}
	}
}

func (s *myServer) HelloError(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	stat := status.New(codes.Unknown, "unknown error")
	stat, _ = stat.WithDetails(&errdetails.DebugInfo{
		Detail: "ここにclient向けのstacktraceが入ってくるイメージ？",
	})
	return &hellopb.HelloResponse{}, stat.Err()
}

func (s *myServer) ChangeOrderPrice(ctx context.Context, req *orderpb.OrderRequest) (*orderpb.OrderResponse, error) {
	return &orderpb.OrderResponse{
		Code:    200,
		Message: fmt.Sprintf("orderId: %d, priceAfterChange: %d, changeReason: %s", req.GetId(), req.GetPriceAfterChange(), req.GetChangeReason()),
	}, nil
}

func (s *myServer) ChangeMultipleOrderPrice(req *orderpb.OrderRequest, stream orderpb.OrderService_ChangeMultipleOrderPriceServer) error {
	resCount := 5

	for i := 0; i < resCount; i++ {
		if err := stream.Send(&orderpb.OrderResponse{
			Code:    200,
			Message: fmt.Sprintf("[%d], orderId: %d, priceAfterChange: %d, changeReason: %s", i, req.GetId(), req.GetPriceAfterChange(), req.GetChangeReason()),
		}); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}

	return nil
}
