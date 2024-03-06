package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage_CreatePoll(t *testing.T) {
	type args struct {
		pollID string
		entity PollEntity
	}
	tests := []struct {
		name string
		s    *MemStorage
		args args
	}{
		{
			name: "Test 1",
			s:    newMemory(),
			args: args{
				pollID: "id1",
				entity: PollEntity{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.s.CreatePoll(tt.args.pollID, tt.args.entity)
			assert.NoError(t, err)

			err = tt.s.CreatePoll(tt.args.pollID, tt.args.entity)
			if assert.Error(t, err) {
				assert.Equal(t, ErrAlreadyExists, err)
			}
		})
	}
}

func TestMemStorage_GetPoll(t *testing.T) {
	type args struct {
		pollID string
	}
	type result struct {
		entity PollEntity
		err    error
	}
	tests := []struct {
		name string
		s    *MemStorage
		args args
		want result
	}{
		{
			name: "Test found",
			s: &MemStorage{
				polls: map[string]PollEntity{
					"id1": {ChatID: 1},
				},
			},
			args: args{
				pollID: "id1",
			},
			want: result{
				entity: PollEntity{ChatID: 1},
				err:    nil,
			},
		},
		{
			name: "Test not found",
			s: &MemStorage{
				polls: map[string]PollEntity{
					"id2": {ChatID: 2},
				},
			},
			args: args{
				pollID: "bad id",
			},
			want: result{
				entity: PollEntity{},
				err:    ErrNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := tt.s.GetPoll(tt.args.pollID)
			if tt.want.err != nil {
				if assert.Error(t, err) {
					assert.Equal(t, tt.want.err, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.entity, entity)
			}
		})
	}
}

func TestMemStorage_DeletePoll(t *testing.T) {
	type args struct {
		pollID string
	}
	tests := []struct {
		name string
		s    *MemStorage
		args args
		want error
	}{
		{
			name: "Test found",
			s: &MemStorage{
				polls: map[string]PollEntity{
					"id1": {ChatID: 1},
				},
			},
			args: args{
				pollID: "id1",
			},
			want: nil,
		},
		{
			name: "Test not found",
			s: &MemStorage{
				polls: map[string]PollEntity{
					"id2": {ChatID: 2},
				},
			},
			args: args{
				pollID: "bad id",
			},
			want: ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.s.DeletePoll(tt.args.pollID)
			if tt.want != nil {
				if assert.Error(t, err) {
					assert.Equal(t, tt.want, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
