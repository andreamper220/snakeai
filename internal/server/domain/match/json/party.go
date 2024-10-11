package json

type PartyJson struct {
	Size      int        `json:"size,omitempty"`
	Width     int        `json:"width"`
	Height    int        `json:"height"`
	Obstacles [][2]int32 `json:"obstacles,omitempty"`
	ToUseById string     `json:"to_use_by_id,omitempty"`
}

type PlayerPartyJson struct {
	PartyId string `json:"party_id"`
}
