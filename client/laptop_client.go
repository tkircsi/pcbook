package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/tkircsi/pcbook/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LaptopClient struct {
	service pb.LaptopServiceClient
}

func NewLaptopClient(cc *grpc.ClientConn) *LaptopClient {
	return &LaptopClient{
		service: pb.NewLaptopServiceClient(cc),
	}
}

func (c *LaptopClient) CreateLaptop(laptop *pb.Laptop) {
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := c.service.CreateLaptop(ctx, req)
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

func (c *LaptopClient) SearchLaptop(filter *pb.Filter) {
	log.Printf("search filter: %v", filter)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.SearchLaptopRequest{
		Filter: filter,
	}
	stream, err := c.service.SearchLaptop(ctx, req)
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

func (c *LaptopClient) UploadImage(laptopID string, imagePath string) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("cannot open image file:", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := c.service.UploadImage(ctx)
	if err != nil {
		log.Fatal("cannot  upload image: ", err)
	}

	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{
				LaptopId:  laptopID,
				ImageType: filepath.Ext(imagePath),
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		log.Fatal("cannot send image info: ", err)
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}
		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send chunk to server: ", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot receive response: ", err)
	}

	log.Printf("image uploaded with id: %s, size: %d", res.GetId(), res.GetSize())
}

func (c *LaptopClient) RateLaptop(laptopIDs []string, scores []float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := c.service.RateLaptop(ctx)
	if err != nil {
		return fmt.Errorf("cannot rate laptop: %v", err)
	}

	waitResponse := make(chan error)
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Println("no more response")
				waitResponse <- nil
				return
			}
			if err != nil {
				waitResponse <- fmt.Errorf("cannot receive stream response: %v", err)
				return
			}

			log.Printf("received response: %v", res)
		}
	}()

	for i, laptopID := range laptopIDs {
		req := &pb.RateLaptopRequest{
			LaptopId: laptopID,
			Score:    scores[i],
		}

		err = stream.Send(req)
		if err != nil {
			return fmt.Errorf("cannot send stream request: %v - %v", err, stream.RecvMsg(nil))
		}
	}

	err = stream.CloseSend()
	if err != nil {
		return fmt.Errorf("cannot close stream: %v", err)
	}

	err = <-waitResponse
	return err
}
