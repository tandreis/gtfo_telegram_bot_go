package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	type args struct {
		storageType string
		storagePath string
	}
	type result struct {
		storage any
		err     error
	}
	tests := []struct {
		name string
		args args
		want result
	}{
		{
			name: "Test memory",
			args: args{storageType: "memory", storagePath: ""},
			want: result{storage: &memStorage{}, err: nil},
		},
		{
			name: "Test bad",
			args: args{storageType: "bad storage", storagePath: ""},
			want: result{storage: nil, err: ErrStorageBad},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage, err := New(tt.args.storageType, tt.args.storagePath)

			require.Equal(t, err, tt.want.err)
			assert.IsType(t, storage, tt.want.storage)
		})
	}
}
