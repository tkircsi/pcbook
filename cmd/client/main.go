package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	"github.com/tkircsi/pcbook/pb"
	"github.com/tkircsi/pcbook/samples"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()

	log.Printf("dial server %s", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	laptopClient := pb.NewLaptopServiceClient(conn)
	for i := 0; i < 10; i++ {
		createLaptop(laptopClient)
	}

	filter := &pb.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz:   2.5,
		MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGYBYTE},
	}

	searchLaptop(laptopClient, filter)
}

func createLaptop(laptopClient pb.LaptopServiceClient) {
	laptop := samples.NewLaptop()
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := laptopClient.CreateLaptop(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Println("laptop already exists")
		} else {
			log.Fatal("cannot create laptop:", err)
		}
		return
	}
	log.Printf("laptop created with id: %s", res.Id)
}

func searchLaptop(laptopClient pb.LaptopServiceClient, filter *pb.Filter) {
	log.Printf("search filter: %v", filter)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.SearchLaptopRequest{
		Filter: filter,
	}
	stream, err := laptopClient.SearchLaptop(ctx, req)
	if err != nil {
		log.Fatal("cannot search laptop: ", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal("cannot receive response: ", err)
		}
		laptop := res.GetLaptop()
		log.Print("- found: ", laptop.GetId())
		log.Print("\t+ brand: ", laptop.GetBrand())
		log.Print("\t+ name: ", laptop.GetName())
		log.Print("\t+ CPU cores: ", laptop.GetCpu().GetNumberCores())
		log.Print("\t+ CPU min GHz: ", laptop.GetCpu().GetMinGhz())
		log.Print("\t+ RAM: ", laptop.GetRam().GetValue(), laptop.GetRam().GetUnit())
		log.Print("\t+ price: ", laptop.GetPriceUsd(), "USD")
	}
}
