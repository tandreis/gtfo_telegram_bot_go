package storage

import "sync"

type MemStorage struct {
	mu    sync.RWMutex
	polls map[string]PollEntity
}

func (s *MemStorage) CreatePoll(pollID string, entity PollEntity) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.polls[pollID]
	if ok {
		return ErrAlreadyExists
	}

	s.polls[pollID] = entity
	return nil
}

func (s *MemStorage) GetPoll(pollID string) (PollEntity, error) {
	s.mu.RLock()
	pollData, ok := s.polls[pollID]
	s.mu.RUnlock()
	if !ok {
		return PollEntity{}, ErrNotFound
	}

	return pollData, nil
}

func (s *MemStorage) DeletePoll(pollID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.polls[pollID]
	delete(s.polls, pollID)

	if !ok {
		return ErrNotFound
	}
	return nil
}

func newMemory() *MemStorage {
	return &MemStorage{polls: make(map[string]PollEntity)}
}
