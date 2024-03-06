package storage

import (
	"errors"
	"time"
)

type PollEntity struct {
	ChatID     int64
	MessageID  int64
	MaxAnswers int
	Started    time.Time
}

type Storage interface {
	CreatePoll(pollID string, entity PollEntity) error
	GetPoll(pollID string) (PollEntity, error)
	DeletePoll(pollID string) error
}

var (
	ErrStorageBad    = errors.New("bad storage type")
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

// New creates new storage of a choosen type
func New(storageType string, storagePath string) (Storage, error) {
	switch storageType {
	case "memory":
		return newMemory(), nil
	default:
		return nil, ErrStorageBad
	}
}
