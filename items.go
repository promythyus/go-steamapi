package steamapi

import (
	"fmt"
	"net/url"
	"strconv"
)

type playerItemsJSON struct {
	Result Inventory
}

// Inventory is the inventory of the user as represented in steam
type Inventory struct {
	Status        uint
	BackpackSlots int `json:"num_backpack_slots"`
	Items         []Item
}

// Item in an inventory
type Item struct {
	ID                uint64
	OriginalID        uint64 `json:"original_id"`
	Defindex          int
	Level             int
	Quantity          int
	Origin            int
	Untradeable       bool   `json:"flag_cannot_trade,omitempty"`
	Uncraftable       bool   `json:"flag_cannot_craft,omitempty"`
	InventoryToken    uint64 `json:"inventory"`
	Quality           int
	CustomName        string      `json:"custom_name,omitempty"`
	CustomDescription string      `json:"custom_description,omitempty"`
	Attributes        []Attribute `json:",omitempty"`
	Equipped          []EquipInfo `json:",omitempty"`
}

// Position gets the position of the item in an inventory
func (i *Item) Position() uint16 {
	return uint16(i.InventoryToken & 0xFFFF)
}

// Attribute is the attribute of an item
type Attribute struct {
	Defindex    int
	Value       int
	FloatValue  float64      `json:"float_value,omitempty"`
	AccountInfo *AccountInfo `json:",omitempty"`
}

// AccountInfo is id and name of user
type AccountInfo struct {
	SteamID     uint64 `json:",string"`
	PersonaName string
}

// EquipInfo class and slot of equipment
type EquipInfo struct {
	Class int
	Slot  int
}

// GetPlayerItems Fetches the player summaries for the given Steam Id.
func GetPlayerItems(baseSteamAPIURL string, steamID uint64, appID uint64, apiKey string) (*Inventory, error) {

	getPlayerItems := NewSteamMethod(baseSteamAPIURL, "IEconItems_"+strconv.FormatUint(appID, 10), "GetPlayerItems", 1)

	vals := url.Values{}
	vals.Add("key", apiKey)
	vals.Add("SteamId", strconv.FormatUint(steamID, 10))

	var resp playerItemsJSON

	err := getPlayerItems.Request(vals, &resp)

	if err != nil {
		return nil, fmt.Errorf("steamapi GetPlayerItems: %v", err)
	}

	return &resp.Result, nil
}
