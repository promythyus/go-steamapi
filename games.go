package steamapi

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type classGameJSON struct {
	Response map[string]json.RawMessage `json:"response"`
}

// Game is the details about a game
type Game struct {
	AppID                    uint   `json:"appid"`
	Name                     string `json:"name"`
	PlaytimeFortnight        uint64 `json:"playtime_2weeks"`
	PlaytimeForever          uint64 `json:"playtime_forever"`
	IconURL                  string `json:"img_icon_url"`
	LogoURL                  string `json:"img_logo_url"`
	HasCommunityVisibleStats bool   `json:"has_community_visible_stats"`
}

// GetOwnedGames returns game details
func GetOwnedGames(steamID uint64, appIDs []uint64, IncludeAppInfo bool, IncludeFreeGames bool, apiKey string) (uint64, []Game, error) {

	var getOwnedGames = NewSteamMethod("IPlayerService", "GetOwnedGames", 1)

	vals := url.Values{}
	vals.Add("key", apiKey)
	vals.Add("steamid", strconv.FormatUint(steamID, 10))

	// Name, IconURL, LogoURL will not be populated when IncludeAppInfo is false
	if IncludeAppInfo {
		vals.Add("include_appinfo", "1")
	}

	if IncludeFreeGames {
		vals.Add("include_played_free_games", "1")
	}

	// Restricts results to the AppIDs specified
	for i, appID := range appIDs {
		vals.Add("appids_filter["+strconv.FormatInt(int64(i), 10)+"]", strconv.FormatUint(appID, 10))
	}

	var resp classGameJSON
	err := getOwnedGames.Request(vals, &resp)
	if err != nil {
		return 0, nil, err
	}

	var games []Game
	var gamesCount uint64

	for index, object := range resp.Response {
		var err error
		switch index {
		case "game_count":
			err = json.Unmarshal(object, &gamesCount)
		case "games":
			err = json.Unmarshal(object, &games)
		}

		if err != nil {
			continue
		}
	}

	return gamesCount, games, nil
}
