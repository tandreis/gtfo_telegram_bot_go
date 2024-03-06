package config

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	confFile := path.Join(dir, "test.yaml")

	err := os.Setenv("CONFIG_PATH", confFile)
	require.NoError(t, err)

	type result struct {
		cfg *Config
		err bool
	}
	tests := []struct {
		name string
		got  string
		want result
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
			want: result{
				cfg: &Config{
					Logger:  Logger{Level: "debug"},
					Storage: Storage{Type: "db", Path: "/tmp/file.db"},
					Bot:     Bot{Token: "token1", MaxPolls: 1},
				},
				err: false,
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
			want: result{
				cfg: &Config{
					Logger:  Logger{Level: "info"},
					Storage: Storage{Type: "memory", Path: ""},
					Bot:     Bot{Token: "token2", MaxPolls: 2},
				},
				err: false,
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
			want: result{
				cfg: &Config{
					Logger:  Logger{Level: "info"},
					Storage: Storage{Type: "memory", Path: ""},
					Bot:     Bot{Token: "token2", MaxPolls: 1},
					Steam:   Steam{ApiKey: "key1"},
				},
				err: false,
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
			want: result{
				cfg: &Config{
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
				err: false,
			},
		},
		{
			name: "Test bad config",
			got:  "",
			want: result{err: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.WriteFile(confFile, []byte(tt.got), 0664)
			require.NoError(t, err)

			got, err := Load()

			if tt.want.err {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, got, tt.want.cfg)
			}

		})
	}
}

func TestLoad_NotFound(t *testing.T) {
	err := os.Setenv("CONFIG_PATH", "!bad file path!")
	require.NoError(t, err)

	_, err = Load()
	assert.Error(t, err)
}

func TestLoad_NoConf(t *testing.T) {
	err := os.Setenv("CONFIG_PATH", "")
	require.NoError(t, err)

	_, err = Load()
	assert.Error(t, err)
}
