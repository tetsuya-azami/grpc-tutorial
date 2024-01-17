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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		fmt.Println("3: send Request to ChangeOrderPrice")
		fmt.Println("4: send Request to ChangeMultipleOrderPrice")
		fmt.Println("5: exit")
		fmt.Println("please enter >")

		scanner.Scan()
		in := scanner.Text()
		switch in {
		case "1":
			Hello()

		case "2":
			HelloServerStream()

		case "3":
			ChangeOrderPrice()

		case "4":
			ChangeMultipleOrderPrice()

		case "5":
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

	res, err := helloClient.Hello(context.Background(), req)
	if err != nil {
		fmt.Println(err)
	} else {
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
