package models

var (
	CreateConn    = 1
	StartGame     = 2
	UpdateCurv    = 3
	EndCurv       = 4
	CollisionDrop = 5
	InvalidDrop   = 6
	DropBall      = 7
	EndGame       = 8
	Close         = 9
	StatusSuccess = "success"
	StatusFailure = "failure"
)

type Ball struct {
	Number int    `json:"number"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	R      int    `json:"r"`
	Type   string `json:"type"`
	Color  string `json:"color"`
}

type Level struct {
	LevelNumber int    `json:"levelNumber"`
	Balls       []Ball `json:"balls"`
}

type Dot struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Message struct {
	Event   int       `json:"event"`
	Players *[]string `json:"players,omitempty"`
	Level   *Level    `json:"level,omitempty"`
	Curve   *[]Dot    `json:"curve,omitempty"`
	Ball    *Ball     `json:"ball,omitempty"`
	Status  *string   `json:"status,omitempty"`
}
