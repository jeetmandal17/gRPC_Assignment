package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gRPC/calc/calc"
	"google.golang.org/grpc"
)



type server struct{
	calc.UnimplementedCalculatorServiceServer
}

// Defininng the server side of each API
// 1. SUM API [UNARY Call]
func (*server) ComputeSum(ctx context.Context,req *calc.SumRequest) (*calc.SumResponse, error){

	fmt.Println("This function was invoked to process sum API")

	// Extracting values from request
	a := req.GetA()
	b := req.GetB()

	// Compute the sum and send back the response
	resp := &calc.SumResponse{
		Sum: a+b,
	}

	return resp, nil
}


// 2. Prime API [Serve-Side Stream]
func (*server) ComputePrime(req *calc.PrimeRequest, resp calc.CalculatorService_ComputePrimeServer) error{
	fmt.Println("This func was invoked to process the Prime API")

	// Extract the value of n from request
	n := req.GetN()

	// Check for all number less than n
	var i int32
	for i=2; i<n; i++{
		if isPrime(i) == true {
			result := &calc.PrimeResponse{
				Number: i,
			}

			// Delay to notice responses
			time.Sleep(1*time.Second)
			resp.Send(result)
		}
	}

	return nil
}

func isPrime(val int32) bool{

	if val <= 1{
		return false
	}

	var i int32
	for i=2; i<val; i++{
		if val % int32(i) == 0{
			return false
		}
	}

	return true
}

func main(){

	// Setup the gRPC server
	fmt.Println("Starting gRPC Service...")

	listen, err := net.Listen("tcp", ":9090")
	if err != nil{
		log.Fatalf("Failed to listen: %v", err)
	}

	// Setup gRPC server
	svr := grpc.NewServer()
	calc.RegisterCalculatorServiceServer(svr, &server{})

	if err = svr.Serve(listen); err != nil{
		log.Fatalf("failed to serve : %v ", err)
	}
}