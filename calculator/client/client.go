package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gRPC/calc/calc"
	"google.golang.org/grpc"
)


func main(){
	cc, err := grpc.Dial(":9090", grpc.WithInsecure())
	if err != nil{
		log.Fatalf("cannot connect to server: %v", err)
	}

	//Close the connection with the server
	defer cc.Close()

	// Register a new client to communicate with the server
	cli := calc.NewCalculatorServiceClient(cc)

	// Request sent to SUM API (Unary)
	//ComputeSum(cli)

	// Request sent to Prime API (Server-side Stream)
	//ComputePrime(cli)

	// Request sent to ComputAvg API (Client-side Stream)
	ComputeAvg(cli)
}

// 1. Sum API
func ComputeSum(cli calc.CalculatorServiceClient){

	fmt.Println("Starting SUM API (Unary)")

	// Form the request to be sent to the server
	req := &calc.SumRequest{
		A: 100,
		B: 200,
	}

	resp, err := cli.ComputeSum(context.Background(), req)
	if err != nil{
		log.Fatalf("error calling SUM gRPC call: %v", err)
	}

	log.Printf("Response from SUM gRPC API: %v", resp.Sum)
}

// 2. Prime API
func ComputePrime(cli calc.CalculatorServiceClient){

	fmt.Println("Starting Prime API (Server-side stream)")

	// Form the request to sent to server
	req := &calc.PrimeRequest{
		N: 10,
	}

	resp, err := cli.ComputePrime(context.Background(), req)
	if err != nil{
		log.Fatalf("cannnot connect to client: %v", err)
	}

	for {
		msg, err := resp.Recv()
		if err == io.EOF{
			break
		}

		if err != nil{
			log.Fatalf("Issue getting data from Prime API: %v", err)
		}

		fmt.Println("Response from gRPC server: ", msg.Number)
	}
}

// 3. AvgCompute API
func ComputeAvg(cli calc.CalculatorServiceClient){

	fmt.Println("Starting Average API (Client-side stream)")

	// Form a connection with the server
	reqStream, err := cli.ComputeAvg(context.Background())
	if err != nil {
		log.Fatal("Error occured while performing client side streaming", err)
	}

	// Form the request to send to server
	req := []*calc.AvgRequest{
		{
			Num: 1,
		},
		{
			Num: 2,
		},
		{
			Num: 3,
		},
		{
			Num: 4,
		},
		{
			Num: 5,
		},
	}

	// Send the items in the reqStream
	for _, item := range req{
		fmt.Println("Sending Requests to Avg gRPC....- > ", item.Num)
		reqStream.Send(item)
		time.Sleep(1*time.Second)
	}

	// Get the final result from the gRPC API
	resp, err := reqStream.CloseAndRecv()
	if err != nil{
		log.Fatal("error while receiving the average : ", err)
	}

	fmt.Println("The average is : ", resp.Result)

}