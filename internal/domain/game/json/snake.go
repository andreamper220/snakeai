package json

type PointJson struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type SnakeJson struct {
	Color string      `json:"color"`
	Body  []PointJson `json:"body"`
}
