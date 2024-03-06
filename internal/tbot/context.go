package tbot

import (
	"context"
	"errors"

	"github.com/tandreis/gtfo_telegram_bot_go/internal/storage"
)

type ctxKeyType string

type ctxData struct {
	storage storage.Storage
}

const ctxKey ctxKeyType = "data"

var (
	errCtxGetDataFailed = errors.New("failed to get context value")
)

func newCtxData(storage storage.Storage) *ctxData {
	return &ctxData{storage: storage}
}

func getStorage(ctx context.Context) (storage.Storage, error) {
	data, ok := ctx.Value(ctxKey).(*ctxData)
	if !ok {
		return nil, errCtxGetDataFailed
	}

	return data.storage, nil
}
