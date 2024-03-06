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
		s    *memStorage
		args args
		want error
	}{
		{
			name: "Test OK",
			s: &memStorage{
				polls: map[string]PollEntity{},
			},
			args: args{
				pollID: "id1",
				entity: PollEntity{},
			},
			want: nil,
		},
		{
			name: "Test exists",
			s: &memStorage{
				polls: map[string]PollEntity{
					"id1": {ChatID: 1},
				},
			},
			args: args{
				pollID: "id1",
				entity: PollEntity{},
			},
			want: ErrAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.s.CreatePoll(tt.args.pollID, tt.args.entity)
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
		s    *memStorage
		args args
		want result
	}{
		{
			name: "Test found",
			s: &memStorage{
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
			s: &memStorage{
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
		s    *memStorage
		args args
		want error
	}{
		{
			name: "Test found",
			s: &memStorage{
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
			s: &memStorage{
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

func Test_memStorage_CreateUser(t *testing.T) {
	type args struct {
		chatID int64
		user   UserEntity
	}
	tests := []struct {
		name string
		s    *memStorage
		args args
		want error
	}{
		{
			name: "Test OK",
			s: &memStorage{
				users: map[int64][]UserEntity{
					1: {
						{Name: "Name1", TelegramID: 11},
						{Name: "Name2", TelegramID: 12},
						{Name: "Name3", TelegramID: 13},
					},
				},
			},
			args: args{
				chatID: 1,
				user:   UserEntity{Name: "Name1", TelegramID: 14},
			},
			want: nil,
		},
		{
			name: "Test empty",
			s: &memStorage{
				users: map[int64][]UserEntity{},
			},
			args: args{
				chatID: 1,
				user:   UserEntity{Name: "Name1", TelegramID: 14},
			},
			want: nil,
		},
		{
			name: "Test exists",
			s: &memStorage{
				users: map[int64][]UserEntity{
					1: {
						{Name: "Name1", TelegramID: 11},
						{Name: "Name2", TelegramID: 12},
						{Name: "Name3", TelegramID: 13},
					},
				},
			},
			args: args{
				chatID: 1,
				user:   UserEntity{Name: "Name1", TelegramID: 12},
			},
			want: ErrAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.s.CreateUser(tt.args.chatID, tt.args.user)
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

func Test_memStorage_GetUsers(t *testing.T) {
	type args struct {
		chatID int64
	}
	type result struct {
		users []UserEntity
		err   error
	}
	tests := []struct {
		name string
		s    *memStorage
		args args
		want result
	}{
		{
			name: "Test found",
			s: &memStorage{
				users: map[int64][]UserEntity{
					1: {
						{Name: "User1"},
						{Name: "User2"},
					},
					2: {{Name: "user3"}},
				},
			},
			args: args{
				chatID: 1,
			},
			want: result{
				users: []UserEntity{
					{Name: "User1"},
					{Name: "User2"},
				},
				err: nil,
			},
		},
		{
			name: "Test not found",
			s: &memStorage{
				users: map[int64][]UserEntity{
					1: {
						{Name: "User1"},
						{Name: "User2"},
					},
					2: {{Name: "user3"}},
				},
			},
			args: args{
				chatID: 3,
			},
			want: result{
				users: []UserEntity{
					{Name: "User1"},
					{Name: "User2"},
				},
				err: ErrNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := tt.s.GetUsers(tt.args.chatID)
			if tt.want.err != nil {
				if assert.Error(t, err) {
					assert.Equal(t, tt.want.err, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.users, entity)
			}
		})
	}
}

func Test_memStorage_DeleteUser(t *testing.T) {
	type args struct {
		chatID     int64
		telegramID int64
	}
	tests := []struct {
		name  string
		s     *memStorage
		args  args
		want  error
		wantS *memStorage
	}{
		{
			name: "Test OK",
			s: &memStorage{
				users: map[int64][]UserEntity{
					1: {
						{Name: "Name1", TelegramID: 11},
						{Name: "Name2", TelegramID: 12},
						{Name: "Name3", TelegramID: 13},
					},
					2: {{Name: "Name4"}},
				},
			},
			wantS: &memStorage{users: map[int64][]UserEntity{
				1: {
					{Name: "Name2", TelegramID: 12},
					{Name: "Name3", TelegramID: 13},
				},
				2: {{Name: "Name4"}},
			}},
			args: args{
				chatID:     1,
				telegramID: 11,
			},
			want: nil,
		},
		{
			name: "Test delete chat",
			s: &memStorage{
				users: map[int64][]UserEntity{
					1: {
						{Name: "Name1", TelegramID: 11},
					},
				},
			},
			wantS: &memStorage{
				users: map[int64][]UserEntity{},
			},
			args: args{
				chatID:     1,
				telegramID: 11,
			},
			want: nil,
		},
		{
			name: "Test not found",
			s: &memStorage{
				users: map[int64][]UserEntity{
					1: {
						{Name: "Name1", TelegramID: 11},
						{Name: "Name2", TelegramID: 12},
						{Name: "Name3", TelegramID: 13},
					},
				},
			},
			wantS: &memStorage{
				users: map[int64][]UserEntity{
					1: {
						{Name: "Name1", TelegramID: 11},
						{Name: "Name2", TelegramID: 12},
						{Name: "Name3", TelegramID: 13},
					},
				},
			},
			args: args{
				chatID:     5,
				telegramID: 11,
			},
			want: ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.s.DeleteUser(tt.args.chatID, tt.args.telegramID)
			if tt.want != nil {
				if assert.Error(t, err) {
					assert.Equal(t, tt.want, err)
				}
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantS, tt.s)
		})
	}
}
