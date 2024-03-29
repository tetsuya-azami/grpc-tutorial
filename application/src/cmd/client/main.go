package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	hellopb "mygrpc/pkg/grpc/hello"
	orderpb "mygrpc/pkg/grpc/order"
	"os"
	"strconv"
	"strings"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	scanner           *bufio.Scanner
	helloClient       hellopb.GreetingServiceClient
	orderChangeClient orderpb.OrderServiceClient
)

func main() {
	fmt.Println("start gRPC client.")

	// 1. 標準入力から文字列を受け取るスキャナを用意
	scanner = bufio.NewScanner(os.Stdin)

	// 2. gRPCサーバーとのコネクションを確立
	address := "localhost:8080"
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)

	if err != nil {
		log.Fatal("Connection failed.")
		return
	}
	defer conn.Close()

	// helloClientを作成
	helloClient = hellopb.NewGreetingServiceClient(conn)
	// orderChangeClientを作成
	orderChangeClient = orderpb.NewOrderServiceClient(conn)

	for {
		fmt.Println("1: send Request to Hello")
		fmt.Println("2: send Request to HelloServerStream")
		fmt.Println("3: send Request to HelloClientStream")
		fmt.Println("4: HelloBiStream")
		fmt.Println("5: HelloError")
		fmt.Println("6: send Request to ChangeOrderPrice")
		fmt.Println("7: send Request to ChangeMultipleOrderPrice")
		fmt.Println("8: exit")
		fmt.Println("please enter >")

		scanner.Scan()
		in := scanner.Text()
		switch in {
		case "1":
			Hello()

		case "2":
			HelloServerStream()

		case "3":
			HelloClientStream()

		case "4":
			HelloBiStream()

		case "5":
			HelloError()

		case "6":
			ChangeOrderPrice()

		case "7":
			ChangeMultipleOrderPrice()

		case "8":
			fmt.Println("bye")
			goto M
		}
	}
M:
}

func Hello() {
	fmt.Println("Please enter your name.")
	scanner.Scan()
	name := scanner.Text()

	req := &hellopb.HelloRequest{
		Name: name,
	}

	ctx := context.Background()
	md := metadata.New(map[string]string{"type": "unary", "from": "client"})
	ctx = metadata.NewOutgoingContext(ctx, md)

	var header, trailer metadata.MD
	res, err := helloClient.Hello(ctx, req, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(header)
		fmt.Println(trailer)
		fmt.Println(res.GetMessage())
		fmt.Println()
	}
}

func HelloServerStream() {
	fmt.Println("Please enter your name.")
	scanner.Scan()
	name := scanner.Text()

	req := &hellopb.HelloRequest{
		Name: name,
	}

	streamClient, err := helloClient.HelloServerStream(context.Background(), req)

	if err != nil {
		return
	}

	for {
		res, err := streamClient.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("all responses has already been received.")
			break
		}

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	}
}

func HelloClientStream() {
	stream, err := helloClient.HelloClientStream(context.Background())

	if err != nil {
		fmt.Println(err)
		return
	}

	resCount := 5

	for i := 0; i < resCount; i++ {
		fmt.Printf("please enter name[%d]\n", i)
		scanner.Scan()
		name := scanner.Text()

		if err := stream.Send(&hellopb.HelloRequest{
			Name: name,
		}); err != nil {
			fmt.Println(err)
			return
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.GetMessage())
	}
}

func HelloBiStream() {
	ctx := context.Background()
	md := metadata.New(map[string]string{"type": "stream", "from": "client"})
	ctx = metadata.NewOutgoingContext(ctx, md)
	stream, err := helloClient.HelloBiStream(ctx)

	if err != nil {
		fmt.Println(err)
		return
	}

	var sendEnd, receiveEnd bool
	sendCount := 0
	for !sendEnd && !receiveEnd {
		// 送信処理
		if !sendEnd {
			fmt.Printf("please enter name[%d]\n", sendCount)
			scanner.Scan()
			name := scanner.Text()

			if err := stream.Send(&hellopb.HelloRequest{
				Name: name,
			}); err != nil {
				fmt.Println(err)
				return
			}

			sendCount++

			if sendCount == 5 {
				sendEnd = true
				if err := stream.CloseSend(); err != nil {
					fmt.Println(err)
					return
				}
			}
		}

		// 受信処理
		if !receiveEnd {
			res, err := stream.Recv()

			if errors.Is(err, io.EOF) {
				fmt.Println("all responses has already been received.")
				sendEnd = true
				break
			}

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(res)
			}
		}
	}
}

func HelloError() {
	fmt.Println("Please enter your name.")
	scanner.Scan()
	name := scanner.Text()

	req := &hellopb.HelloRequest{
		Name: name,
	}

	res, err := helloClient.HelloError(context.Background(), req)

	if err != nil {
		if stat, ok := status.FromError(err); ok {
			fmt.Printf("Code: %s\n", stat.Code())
			fmt.Printf("Message: %s\n", stat.Message())

			for _, detail := range stat.Details() {
				switch t := detail.(type) {
				case *errdetails.BadRequest:
					fmt.Printf("BadRequest: %s\n", t)
				case *errdetails.DebugInfo:
					fmt.Printf("DebugInfo: %s\n", t)
				case *errdetails.RetryInfo:
					fmt.Printf("RetryInfo: %s\n", t)
				default:
					fmt.Printf("Unexpected type: %s\n", t)
				}
			}
		} else {
			fmt.Println(err)
		}
	}

	fmt.Println(res.GetMessage())
	fmt.Println()
}

func ChangeOrderPrice() {
	fmt.Println("Please enter order id.")
	scanner.Scan()
	idUint64 := ParseUint(scanner.Text())

	fmt.Println("Please enter order priceAfterChange.")
	scanner.Scan()
	priceAfterChange := ParseUint(scanner.Text())

	fmt.Println("Please enter changeReason.")
	scanner.Scan()
	changeReason, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		fmt.Println(err)
	}

	req := &orderpb.OrderRequest{
		Id:               idUint64,
		PriceAfterChange: priceAfterChange,
		ChangeReason:     orderpb.OrderChangeReason(changeReason),
	}

	res, err := orderChangeClient.ChangeOrderPrice(context.Background(), req)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("code: %d, message: %s", res.GetCode(), res.GetMessage())
		fmt.Println()
	}
}

func ChangeMultipleOrderPrice() {
	fmt.Println("Please enter order id.")
	scanner.Scan()
	id := ParseUint(scanner.Text())

	fmt.Println("Please enter order priceAfterChange.")
	scanner.Scan()
	priceAfterChange := ParseUint(scanner.Text())

	fmt.Println("Please enter order changeReason.")
	scanner.Scan()
	changeReason := scanner.Text()

	stream, err := orderChangeClient.ChangeMultipleOrderPrice(context.Background(), &orderpb.OrderRequest{
		Id:               id,
		PriceAfterChange: priceAfterChange,
		ChangeReason:     orderpb.OrderChangeReason(ParseUint(changeReason)),
	})

	if err != nil {
		fmt.Println(err)
	}

	for {
		res, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("all responses has already been received.")
			break
		}

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	}
}

func ParseUint(source string) uint64 {
	idUint64, err := strconv.ParseUint(source, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	return idUint64
}
