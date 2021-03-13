package service

import "sync"

// RatingStore interface defines the rating storage methods
type RatingStore interface {
	// Add adds a new rating score to the laptop
	Add(laptop_id string, score float64) (*Rating, error)
}

// Rating contains the rating information of a laptop
type Rating struct {
	Count uint32
	Sum   float64
}

// InMemoryRatingStore in memory implementation of a
// RatingStore interface
type InMemoryRatingStore struct {
	mutex  sync.RWMutex
	rating map[string]*Rating
}

func NewInMemoryRatingStore() *InMemoryRatingStore {
	return &InMemoryRatingStore{
		rating: make(map[string]*Rating),
	}
}

func (s *InMemoryRatingStore) Add(laptop_id string, score float64) (*Rating, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	rating := s.rating[laptop_id]
	if rating == nil {
		rating = &Rating{
			Count: 1,
			Sum:   score,
		}
	} else {
		rating.Count++
		rating.Sum += score
	}

	s.rating[laptop_id] = rating
	return rating, nil
}
