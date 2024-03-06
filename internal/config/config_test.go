package config

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMustLoad(t *testing.T) {
	dir := t.TempDir()
	confFile := path.Join(dir, "test.yaml")

	err := os.Setenv("CONFIG_PATH", confFile)
	require.NoError(t, err)

	tests := []struct {
		name string
		got  string
		want *Config
	}{
		{
			name: "Test 1",
			got: `
logger:
  level: "debug"
storage:
  type: "db"
  path: "/tmp/file.db"
bot:
  token: "token1"
`,
			want: &Config{
				Logger:  Logger{Level: "debug"},
				Storage: Storage{Type: "db", Path: "/tmp/file.db"},
				Bot:     Bot{Token: "token1", MaxPolls: 1},
			},
		},
		{
			name: "Test 1",
			got: `
storage:
  type: "memory"
bot:
  token: "token2"
  max_polls: 2
`,
			want: &Config{
				Logger:  Logger{Level: "info"},
				Storage: Storage{Type: "memory", Path: ""},
				Bot:     Bot{Token: "token2", MaxPolls: 2},
			},
		},
		{
			name: "Test 1",
			got: `
storage:
  type: "memory"
bot:
  token: "token2"
steam:
  api_key: "key1"
`,
			want: &Config{
				Logger:  Logger{Level: "info"},
				Storage: Storage{Type: "memory", Path: ""},
				Bot:     Bot{Token: "token2", MaxPolls: 1},
				Steam:   Steam{ApiKey: "key1"},
			},
		},
		{
			name: "Test 1",
			got: `
storage:
  type: "memory"
bot:
  token: "token2"
steam:
  api_key: "key2"
  users:
    - name: "user1"
      steam_id: "sid1"
      telegram_id: 1
      chat_id: 123
`,
			want: &Config{
				Logger:  Logger{Level: "info"},
				Storage: Storage{Type: "memory", Path: ""},
				Bot:     Bot{Token: "token2", MaxPolls: 1},
				Steam: Steam{
					ApiKey: "key2",
					Users: []User{{
						Name:       "user1",
						SteamID:    "sid1",
						TelegramID: 1,
						ChatID:     123,
					}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.WriteFile(confFile, []byte(tt.got), 0664)
			require.NoError(t, err)

			got := MustLoad()
			assert.Equal(t, got, tt.want)
		})
	}
}
