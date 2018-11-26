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
	GetLevel      = 10
	StatusSuccess = "success"
	StatusFailure = "failure"
)

type Circle struct {
	Number int    `json:"num"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	R      int    `json:"r"`
	Type   string `json:"type"`
	Color  string `json:"color"`
}

type Level struct {
	LevelNumber int      `json:"levelNumber"`
	Circles     []Circle `json:"circles,omitempty"`
}

type Dot struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Message struct {
	Event   int       `json:"event"`
	Players *[]string `json:"players,omitempty"`
	Level   *Level    `json:"level,omitempty"`
	Line    *[]Dot    `json:"line,omitempty"`
	Circle  *Circle   `json:"circle,omitempty"`
	Status  *string   `json:"status,omitempty"`
}
