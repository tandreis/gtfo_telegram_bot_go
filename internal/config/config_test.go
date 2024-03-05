package config

import (
	"os"
	"path"
	"reflect"
	"testing"
)

func TestMustLoad(t *testing.T) {
	dir := t.TempDir()
	confFile := path.Join(dir, "test.yaml")
	os.Setenv("CONFIG_PATH", confFile)

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
bot:
  token: "token1"
`,
			want: &Config{
				Logger: Logger{Level: "debug"},
				Bot:    Bot{Token: "token1", MaxPolls: 1},
			},
		},
		{
			name: "Test 1",
			got: `
bot:
  token: "token2"
  max_polls: 2
`,
			want: &Config{
				Logger: Logger{Level: "info"},
				Bot:    Bot{Token: "token2", MaxPolls: 2},
			},
		},
		{
			name: "Test 1",
			got: `
bot:
  token: "token2"
steam:
  api_key: "key1"
`,
			want: &Config{
				Logger: Logger{Level: "info"},
				Bot:    Bot{Token: "token2", MaxPolls: 1},
				Steam:  Steam{ApiKey: "key1"},
			},
		},
		{
			name: "Test 1",
			got: `
bot:
  token: "token2"
steam:
  api_key: "key2"
  users:
    - name: "user1"
      sid: "sid1"
      tid: 1
`,
			want: &Config{
				Logger: Logger{Level: "info"},
				Bot:    Bot{Token: "token2", MaxPolls: 1},
				Steam: Steam{
					ApiKey: "key2",
					Users: []User{{
						Name:       "user1",
						SteamID:    "sid1",
						TelegramID: 1,
					}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.WriteFile(confFile, []byte(tt.got), 0664)
			if got := MustLoad(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MustLoad() = %v, want %v", got, tt.want)
			}
		})
	}
}
