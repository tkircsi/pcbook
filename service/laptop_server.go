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

type LaptopServer struct {
	Store LaptopStore
}

func NewLaptopServer() *LaptopServer {
	return &LaptopServer{}
}

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
