package tbot

import (
	"context"
	"errors"

	"github.com/tandreis/gtfo_telegram_bot_go/internal/storage"
)

type ctxKeyType string

type ctxData struct {
	storage    storage.Storage
	steamToken string
}

const ctxKey ctxKeyType = "data"

var (
	errCtxGetDataFailed = errors.New("failed to get context value")
)

func newCtxData(storage storage.Storage, token string) *ctxData {
	return &ctxData{storage: storage, steamToken: token}
}

func getStorage(ctx context.Context) (storage.Storage, error) {
	data, ok := ctx.Value(ctxKey).(*ctxData)
	if !ok {
		return nil, errCtxGetDataFailed
	}

	return data.storage, nil
}

func getSteamToken(ctx context.Context) (string, error) {
	data, ok := ctx.Value(ctxKey).(*ctxData)
	if !ok {
		return "", errCtxGetDataFailed
	}

	return data.steamToken, nil
}
