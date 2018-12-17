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
	DataLoaded    = 11
	StartInput    = 12
	FinishInput   = 13
	LineInputted  = 14
	GameProcess   = 15
	GameFinish    = 16
	PlayerSuccess = 17
	PlayerFailure = 18
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

type BaseLine struct {
	Dots    []Dot `json:"points"`
	BaseDot Dot   `json:"base_point"`
}

type Message struct {
	Event   int       `json:"event"`
	Players *[]string `json:"players,omitempty"`
	Level   *Level    `json:"level,omitempty"`
	Line    *BaseLine `json:"line,omitempty"`
	Circle  *Circle   `json:"circle,omitempty"`
	Status  *string   `json:"status,omitempty"`
}

type LevelInStorage struct {
	Number int    `db:"number"`
	Level  string `db:"level"`
}

type PlayerLogic struct {
	Line      *BaseLine
	IsReady   bool
	IsFailure bool
}
