package service

import "sync"

// UserStore defines the store methods
type UserStore interface {
	// Save saves a user into the store
	Save(user *User) error
	// Find search a user with the given username in the store
	Find(username string) (*User, error)
}

// InMemoryUserStore soters users in memory
type InMemoryUserStore struct {
	mutex sync.RWMutex
	users map[string]*User
}

// NewInMemoryUserStore creates a InMemoryUser store
func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: make(map[string]*User),
	}
}

func (s *InMemoryUserStore) Save(user *User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.users[user.Username]; ok {
		return ErrAlreadyExists
	}

	s.users[user.Username] = user.Clone()
	return nil
}

func (s *InMemoryUserStore) Find(username string) (*User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if user, ok := s.users[username]; !ok {
		return nil, nil
	} else {
		return user.Clone(), nil
	}
}
