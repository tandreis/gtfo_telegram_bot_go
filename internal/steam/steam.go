package steam

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	steamURL   = "https://api.steampowered.com"
	summaryAPI = "ISteamUser/GetPlayerSummaries/v0002"
)

var states = [...]string{"Offline", "Online", "Busy", "Away"}

type SteamUsers struct {
	Response struct {
		Players []struct {
			Name       string `json:"personaname"`
			GameName   string `json:"gameextrainfo"`
			GameID     string `json:"gameid"`
			LastLogoff int64  `json:"lastlogoff"`
			State      int    `json:"personastate"`
			StateStr   string
			GameURL    string
			Online     bool
		} `json:"players"`
	} `json:"response"`
}

func getJson(url string, target any) error {
	var client = &http.Client{Timeout: 2 * time.Second}
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

// GetPlayerSummaries returns sumary info for a given Steam users
func GetPlayerSummaries(apiKey string, steamIDs []string) SteamUsers {
	IDs := strings.Join(steamIDs, ",")
	url := fmt.Sprintf("%s/%s?key=%s&steamids=[%s]",
		steamURL, summaryAPI, apiKey, IDs)

	var users SteamUsers
	getJson(url, &users)

	for i := range len(users.Response.Players) {
		if users.Response.Players[i].GameName != "" {
			users.Response.Players[i].GameURL =
				fmt.Sprintf("https://store.steampowered.com/app/%s/",
					users.Response.Players[i].GameID)
		}
		users.Response.Players[i].StateStr =
			states[users.Response.Players[i].State]
		users.Response.Players[i].Online = users.Response.Players[i].State != 0
	}

	return users
}
