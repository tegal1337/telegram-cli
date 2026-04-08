package store

import (
	"sync"

	"github.com/zelenin/go-tdlib/client"
)

// UserStore caches user information.
type UserStore struct {
	mu    sync.RWMutex
	users map[int64]*client.User
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[int64]*client.User),
	}
}

// Set adds or updates a user in the cache.
func (s *UserStore) Set(user *client.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[user.ID] = user
}

// Get returns a cached user by ID.
func (s *UserStore) Get(userID int64) (*client.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[userID]
	return user, ok
}

// DisplayName returns a formatted display name for a user.
func (s *UserStore) DisplayName(userID int64) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[userID]
	if !ok {
		return "Unknown"
	}

	name := user.FirstName
	if user.LastName != "" {
		name += " " + user.LastName
	}
	return name
}

// IsOnline checks if a user is currently online.
func (s *UserStore) IsOnline(userID int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[userID]
	if !ok {
		return false
	}

	_, online := user.Status.(*client.UserStatusOnline)
	return online
}

// UpdateStatus updates a user's online status.
func (s *UserStore) UpdateStatus(userID int64, status client.UserStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if user, ok := s.users[userID]; ok {
		user.Status = status
	}
}

// All returns all cached users.
func (s *UserStore) All() []*client.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*client.User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u)
	}
	return users
}
