package service

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/tkircsi/pcbook/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// The LaptopServer holds storage for Laptops and serves the methods
type LaptopServer struct {
	Store LaptopStore
	pb.UnimplementedLaptopServiceServer
}

// NewLaptopServer creates a new server objecz
func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{
		Store: store,
	}
}

// CreateLaptop creates and save a Laptop object into the server's store
func (s *LaptopServer) CreateLaptop(
	ctx context.Context,
	r *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {
	laptop := r.GetLaptop()
	log.Printf("receive a create-laptop request with id: %s", laptop.GetId())

	if len(laptop.GetId()) > 0 {
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid laptop ID: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate new laptop ID: %v", err)
		}
		laptop.Id = id.String()
	}

	// pretend a long call
	// time.Sleep(6 * time.Second)

	if ctx.Err() == context.Canceled {
		log.Println("request is cancelled")
		return nil, status.Error(codes.Canceled, "request is cancelled")
	}

	if ctx.Err() == context.DeadlineExceeded {
		log.Println("deadline is exceeded")
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	err := s.Store.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannot save laptop to the store: %v", err)
	}

	log.Printf("saved laptop with id: %s", laptop.Id)

	res := &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}

	return res, nil
}

// SearchLaptop is server-streaming RPC to search for laptops
func (s *LaptopServer) SearchLaptop(
	r *pb.SearchLaptopRequest,
	stream pb.LaptopService_SearchLaptopServer) error {
	filter := r.GetFilter()
	log.Printf("receive a search-laptop request with filter: %v", filter)

	err := s.Store.Search(
		stream.Context(),
		filter,
		func(laptop *pb.Laptop) error {
			res := &pb.SearchLaptopResponse{
				Laptop: laptop,
			}
			err := stream.Send(res)
			if err != nil {
				return err
			}
			log.Printf("sent laptop with id: %s", laptop.GetId())
			return nil
		},
	)

	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}
	return nil
}
