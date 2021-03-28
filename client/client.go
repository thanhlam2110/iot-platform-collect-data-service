package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	calculatorpb "github.com/thanhlam/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	//cc, err := grpc.Dial("13.212.172.127:50069", grpc.WithInsecure())
	cc, err := grpc.Dial("0.0.0.0:50069", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Err while dial $v", err)
	}
	defer cc.Close()
	client := calculatorpb.NewCalculatorServiceClient(cc)
	//log.Printf("Service client %f", client)
	//unary
	//callSum(client)
	//callPND(client)
	callAverage(client)
	//callFindMax(client)
}
func callSum(c calculatorpb.CalculatorServiceClient) {
	log.Println("Call sum API")
	resp, err := c.Sum(context.Background(), &calculatorpb.SumRequest{
		Num1: 5,
		Num2: 6,
	})
	if err != nil {
		log.Fatal("Call sum wrong %v", err)
	}
	log.Println("Call sum API response", resp.GetResult())
}
func callPND(c calculatorpb.CalculatorServiceClient) {
	log.Println("calling pnd api")
	stream, err := c.PrimeNumberDecomposition(context.Background(), &calculatorpb.PNDRequest{
		Number: 120,
	})

	if err != nil {
		log.Fatalf("callPND err %v", err)
	}

	for {
		resp, recvErr := stream.Recv()
		if recvErr == io.EOF {
			log.Println("server finish streaming")
			return
		}

		if recvErr != nil {
			log.Fatalf("callPND recvErr %v", recvErr)
		}

		log.Printf("prime number %v", resp.GetNumber())
	}
}
func callAverage(c calculatorpb.CalculatorServiceClient) {
	log.Println("Call average API")
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatal("Call average err %v", err)
	}
	/*listReq := []calculatorpb.AverageRequest{
		calculatorpb.AverageRequest{
			Num: 5,
		},
		calculatorpb.AverageRequest{
			Num: 10,
		},
		calculatorpb.AverageRequest{
			Num: 12,
		},
		calculatorpb.AverageRequest{
			Num: 3,
		},
		calculatorpb.AverageRequest{
			Num: 4.2,
		},
	}*/
	listReq := []calculatorpb.AverageRequest{}
	for i := 0; i < 1000000; i++ {
		j := calculatorpb.AverageRequest{
			Num: float32(i),
		}
		listReq = append(listReq, j)
	}
	//fmt.Printf("%T", listReq)
	//gửi đi
	fmt.Println(time.Now())
	for _, req := range listReq {
		err := stream.Send(&req)
		if err != nil {
			log.Fatal("Send average request err %v", err)
		}
	}
	//gửi hết thì nhận response
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("Receive average response err %v", err)
	}
	fmt.Printf("Average response %v", resp)

}
func callFindMax(c calculatorpb.CalculatorServiceClient) {
	log.Println("Call Find Max......")
	stream, err := c.FindMax(context.Background())
	if err != nil {
		log.Fatal("Call Find Max err %v", err)
	}
	//xu ly de dung ham go func()
	waitc := make(chan struct{})

	//step 1: tao request
	/*listReq := []calculatorpb.FindMaxRequest{
		calculatorpb.FindMaxRequest{
			Num: 5,
		},
		calculatorpb.FindMaxRequest{
			Num: 10,
		},
		calculatorpb.FindMaxRequest{
			Num: 12,
		},
		calculatorpb.FindMaxRequest{
			Num: 3,
		},
		calculatorpb.FindMaxRequest{
			Num: 4,
		},
	}*/
	go func() {
		// xử ly gửi nhieu request
		//step 1: tao request
		listReq := []calculatorpb.FindMaxRequest{
			calculatorpb.FindMaxRequest{
				Num: 5,
			},
			calculatorpb.FindMaxRequest{
				Num: 10,
			},
			calculatorpb.FindMaxRequest{
				Num: 12,
			},
			calculatorpb.FindMaxRequest{
				Num: 3,
			},
			calculatorpb.FindMaxRequest{
				Num: 4,
			},
		}
		//Step2: gửi đi
		for _, req := range listReq {
			err := stream.Send(&req)
			if err != nil {
				log.Fatal("Send find max request err %v", err)
				break
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend() //báo cho server đã gửi xong
	}()
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				//server hoàn thành response
				log.Println("Ending findmax api")
				break
			}
			if err != nil {
				log.Fatalf("Receive FindMax err %v", err)
				break
			}
			log.Printf("Max: %v", resp.GetMax())
		}
		close(waitc)
	}()
	<-waitc

}
