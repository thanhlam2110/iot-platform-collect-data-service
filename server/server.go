package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	calculatorpb "github.com/thanhlam/calculator/calculatorpb"
)

type server struct{}

//Unary
func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("Sum Call")
	resp := &calculatorpb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}
	return resp, nil
}

//Server stream
func (*server) PrimeNumberDecomposition(req *calculatorpb.PNDRequest,
	stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	log.Println("PrimeNumberDecomposition called...")
	k := int32(2)
	N := req.GetNumber()
	for N > 1 {
		if N%k == 0 {
			N = N / k
			//send to client
			stream.Send(&calculatorpb.PNDResponse{
				Number: k,
			})
		} else {
			k++
			log.Printf("k increase to %v", k)
		}
	}
	return nil
}
func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	log.Println("Average called...")
	var total float32
	var count int
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// tinh trung binh va return cho client
			resp := &calculatorpb.AverageResponse{
				Result: total / float32(count),
			}
			fmt.Println(time.Now())
			return stream.SendAndClose(resp)
		}
		if err != nil {
			log.Fatal("Err while Recv Average %v", err)
			return err
		}
		log.Println("Receive num %v", req)
		total += req.GetNum()
		count++
	}

}
func (*server) FindMax(stream calculatorpb.CalculatorService_FindMaxServer) error {
	log.Println("Call API FindMax...")
	max := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("EOF...")
			return nil //client ko gửi nữa thì dừng lun ko gửi response
		}
		if err != nil {
			log.Fatalf("Err while Recv Max %v", err)
			return err
		}
		num := req.GetNum()
		log.Printf("Recv num %v\n", num)
		if num > max {
			max = num
		}
		err = stream.Send(&calculatorpb.FindMaxResponse{
			Max: max,
		})
		if err != nil {
			log.Fatalf("Send Max Error %v", err)
			return err
		}
	}
}
func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50069")
	if err != nil {
		log.Fatal("Error while create listen %v", err)
	}
	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{}) // ham nay tu protoc go tao ra

	fmt.Println("Calculator is running...")

	err = s.Serve(lis) //run server tren port

	if err != nil {
		log.Fatal("Error while server %v", err)
	}
}
