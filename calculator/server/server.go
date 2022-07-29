package main

import (
	"context"
	"fmt"
	"io"
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
	fmt.Println("This function was invoked to process the Prime API")

	// Extract the value of n from request
	n := req.GetN()

	// Check for all number less than n
	var i int32
	for i=2; i<n; i++{
		if isPrime(i) {
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

// 3. AvgCompute [Client-side Stream]
func (*server) ComputeAvg(req calc.CalculatorService_ComputeAvgServer) error {

	fmt.Println("This function was invoked to process the ComputeAvg API")

	// Recieve the stream of data from the client
	var finalSum int32 = 0
	var finalCount int32 = 0

	for {
		reqPacket, err := req.Recv()
		if err == io.EOF{
			break
		}

		if err != nil{
			log.Fatal("Error while fetching requests from client", err)
		}

		finalCount++
		finalSum += reqPacket.GetNum()
	}
	
	// compute the average
	avgResult := &calc.AvgResponse{
		Result: (finalSum/finalCount),
	}
	// Send the response and close the channel
	err := req.SendAndClose(avgResult)
	if err != nil{
		log.Fatal("error while sending data and closing connection", err)
	}

	return nil
}

// 4. FMN [BiDi Stream]
func (*server)ComputeFMN(req calc.CalculatorService_ComputeFMNServer) error {

	fmt.Println("This function was invoked to process the FMN API")

	var maxi int32 = -1
	// Recieve the requests from the client
		
	for {
		reqPack, err := req.Recv()
		if err == io.EOF{
			break
		}

		if err != nil{
			log.Fatal("error while getting requests from client : ", err)
		}
	
		// Update the maxi according to the request
		fmt.Println("Received the new Value : ", reqPack.GetNum())
		if maxi < reqPack.GetNum() {
			maxi = reqPack.GetNum()

			// Whenever there is an update, send the update back to client
			newResp := &calc.FMNResponse{
				NewMax: maxi,
			}
		
			err = req.Send(newResp)
			if err != nil {
				log.Fatal("error while sending newMax to client", err)
			}
		}
	}
	return nil
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