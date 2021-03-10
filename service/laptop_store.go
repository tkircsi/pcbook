package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/jinzhu/copier"
	"github.com/tkircsi/pcbook/pb"
)

var (
	ErrAlreadyExists = errors.New("record already exists")
)

type LaptopStore interface {
	Save(*pb.Laptop) error
}

type InMemoryLaptopStore struct {
	mutex sync.RWMutex
	data  map[string]*pb.Laptop
}

func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*pb.Laptop),
	}
}

func (ls *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()

	if _, ok := ls.data[laptop.Id]; ok {
		return ErrAlreadyExists
	}

	// deep copy copier package
	newLaptop := &pb.Laptop{}
	err := copier.Copy(newLaptop, laptop)
	if err != nil {
		return fmt.Errorf("failed to copy laptop data: %v", err)
	}
	ls.data[newLaptop.Id] = newLaptop
	return nil
}
