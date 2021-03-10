package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/jinzhu/copier"
	"github.com/tkircsi/pcbook/pb"
)

var (
	// ErrAlreadyExists Custom error if the Laptop already exists in the store
	ErrAlreadyExists = errors.New("record already exists")
)

// LaptopStore stores the Laptop objects
type LaptopStore interface {
	Save(*pb.Laptop) error
	Find(id string) (*pb.Laptop, error)
}

// InMemoryLaptopStore is an in-memory store implementation
type InMemoryLaptopStore struct {
	mutex sync.RWMutex
	data  map[string]*pb.Laptop
}

// NewInMemoryLaptopStore creates a new InMemoryLaptopStore
func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*pb.Laptop),
	}
}

// Save saves Laptop object into the store
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

// Find search for a laptop in the store by ID
func (ls *InMemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	ls.mutex.RLock()
	defer ls.mutex.RUnlock()

	laptop, ok := ls.data[id]
	if !ok {
		return nil, nil
	}
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data: %w", err)
	}
	return other, nil

}
