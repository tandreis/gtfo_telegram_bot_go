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

type UserEntity struct {
	Name       string
	SteamID    string
	TelegramID int64
}

type Storage interface {
	// CreatePoll adds poll record to storage.
	// Returs error if poll already exists.
	CreatePoll(pollID string, entity PollEntity) error

	// GetPoll reads poll record from storage.
	// Returs error if poll was not found.
	GetPoll(pollID string) (PollEntity, error)

	// DeletePoll removes poll record from storage.
	// Returs error if poll was not found (nothing to delete).
	DeletePoll(pollID string) error

	// CreateUser adds user record to storage.
	// Returs error if user already exists in chat.
	CreateUser(chatID int64, user UserEntity) error

	// GetUsers reads user records from storage for a given chatID.
	// Returs error if chat was not found.
	GetUsers(chatID int64) ([]UserEntity, error)

	// DeleteUser removes user record from storage for a givet chatID.
	// Returs error if user was not found (nothing to delete).
	DeleteUser(chatID int64, telegramID int64) error
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
