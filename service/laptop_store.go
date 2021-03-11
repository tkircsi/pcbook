package service

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	// Save saves the laptop to the store
	Save(*pb.Laptop) error
	// Find finds a laptop by ID
	Find(id string) (*pb.Laptop, error)
	// Search searches for laptop with filter and returns one by one via the found function
	Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error
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
	newLaptop, err := deepCopy(laptop)
	if err != nil {
		return err
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
	return deepCopy(laptop)
}

// Search searches for laptop with filter and returns one by one via the found function
func (ls *InMemoryLaptopStore) Search(
	ctx context.Context,
	filter *pb.Filter,
	found func(laptop *pb.Laptop) error) error {
	ls.mutex.RLock()
	defer ls.mutex.RUnlock()

	for _, laptop := range ls.data {
		// simulate heavy processing
		// time.Sleep(time.Second)
		log.Println("checking laptop id: ", laptop.GetId())

		if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
			log.Println("context is cancelled")
			return errors.New("context is cancelled")
		}

		if isQualified(filter, laptop) {
			// deepCopy
			other, err := deepCopy(laptop)
			if err != nil {
				return err
			}
			err = found(other)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isQualified(filter *pb.Filter, laptop *pb.Laptop) bool {
	if laptop.GetPriceUsd() > filter.GetMaxPriceUsd() {
		return false
	}

	if laptop.GetCpu().GetNumberCores() < filter.GetMinCpuCores() {
		return false
	}

	if laptop.GetCpu().GetMinGhz() < filter.GetMinCpuGhz() {
		return false
	}

	if toBit(laptop.GetRam()) < toBit(filter.GetMinRam()) {
		return false
	}
	return true
}

func toBit(mem *pb.Memory) uint64 {
	value := mem.GetValue()

	switch mem.GetUnit() {
	case pb.Memory_BIT:
		return value
	case pb.Memory_BYTE:
		return value << 3 // value * 8
	case pb.Memory_KILOBYTE:
		return value << 13 // value * 8 * 1024
	case pb.Memory_MEGABYTE:
		return value << 23
	case pb.Memory_GIGYBYTE:
		return value << 33
	case pb.Memory_TERABYTE:
		return value << 43
	default:
		return 0
	}
}

func deepCopy(laptop *pb.Laptop) (*pb.Laptop, error) {
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data: %w", err)
	}
	return other, nil
}
