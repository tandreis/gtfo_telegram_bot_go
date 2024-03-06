package storage

import (
	"slices"
	"sync"
)

type memStorage struct {
	muPolls sync.RWMutex
	muUsers sync.RWMutex
	polls   map[string]PollEntity
	users   map[int64][]UserEntity
}

func (s *memStorage) CreatePoll(pollID string, entity PollEntity) error {
	s.muPolls.Lock()
	defer s.muPolls.Unlock()

	_, ok := s.polls[pollID]
	if ok {
		return ErrAlreadyExists
	}

	s.polls[pollID] = entity
	return nil
}

func (s *memStorage) GetPoll(pollID string) (PollEntity, error) {
	s.muPolls.RLock()
	pollData, ok := s.polls[pollID]
	s.muPolls.RUnlock()
	if !ok {
		return PollEntity{}, ErrNotFound
	}

	return pollData, nil
}

func (s *memStorage) DeletePoll(pollID string) error {
	s.muPolls.Lock()
	defer s.muPolls.Unlock()

	_, ok := s.polls[pollID]
	if !ok {
		return ErrNotFound
	}

	delete(s.polls, pollID)

	return nil
}

func (s *memStorage) CreateUser(chatID int64, user UserEntity) error {
	s.muUsers.Lock()
	defer s.muUsers.Unlock()

	_, ok := s.users[chatID]
	if !ok {
		s.users[chatID] = make([]UserEntity, 0, 10)
	} else {
		if idx := slices.IndexFunc(s.users[chatID],
			func(u UserEntity) bool { return u.TelegramID == user.TelegramID },
		); idx != -1 {
			return ErrAlreadyExists
		}
	}

	s.users[chatID] = append(s.users[chatID], user)
	return nil
}

func (s *memStorage) GetUsers(chatID int64) ([]UserEntity, error) {
	s.muUsers.RLock()
	users, ok := s.users[chatID]
	s.muUsers.RUnlock()
	if !ok {
		return []UserEntity{}, ErrNotFound
	}

	return users, nil
}

func (s *memStorage) DeleteUser(chatID int64, telegramID int64) error {
	s.muUsers.Lock()
	defer s.muUsers.Unlock()

	users, ok := s.users[chatID]
	if !ok {
		return ErrNotFound
	}

	if idx := slices.IndexFunc(s.users[chatID],
		func(u UserEntity) bool { return u.TelegramID == telegramID },
	); idx != -1 {
		copy(users[idx:], users[idx+1:])
		users = users[:len(users)-1]
	}

	if len(users) == 0 {
		delete(s.users, chatID)
	} else {
		s.users[chatID] = users
	}

	return nil
}

func newMemory() *memStorage {
	return &memStorage{
		polls: make(map[string]PollEntity),
		users: make(map[int64][]UserEntity),
	}
}
