package json

type PartyJson struct {
	Size      int        `json:"size,omitempty"`
	Width     int        `json:"width"`
	Height    int        `json:"height"`
	Obstacles [][2]int32 `json:"obstacles,omitempty"`
}
