package models

import (
	"errors"
	"log"
)

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

const (
	WINDOW_W               = 1280
	WINDOW_H               = 720
	MAX_LINE_POINTS_LENGTH = 1500
)

type Window struct {
	width  int
	height int
}

var window Window = Window{WINDOW_W, WINDOW_H}

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

type Dot struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func Sum(d1 Dot, d2 Dot) Dot {
	return Dot{d1.X + d2.X, d1.Y + d2.Y}
}

type BaseLine struct {
	Dots            []Dot `json:"points"`
	BaseDot         Dot   `json:"base_point"`
	CurrentPosition int
}

func (bl BaseLine) Size() int {
	return len(bl.Dots)
}

func (bl BaseLine) CopyDots() []Dot {
	dots := make([]Dot, bl.Size())
	copy(bl.Dots, dots)
	return dots
}

func (bl BaseLine) Copy(dest *BaseLine) {
	dest.CurrentPosition = bl.CurrentPosition
	dest.BaseDot = bl.BaseDot
	dest.Dots = bl.CopyDots()
}

func (bl BaseLine) GetAbsoluteDot(index int) (Dot, error) {
	if index < 0 || index > len(bl.Dots) {
		return Dot{}, errors.New("out of index")
	}

	absDot := Sum(bl.Dots[index], bl.BaseDot)
	return absDot, nil
}

type Line struct {
	BeginLine  BaseLine
	EndLine    BaseLine
	OriginLine BaseLine
	IsReversed bool
}

func (l *Line) Step() bool {
	l.BeginLine.CurrentPosition++
	l.EndLine.CurrentPosition++

	if l.BeginLine.CurrentPosition == l.BeginLine.Size() {
		l.BeginLine.Copy(&l.EndLine)
		//construct endLine
	}

	return !l.isLineOutOfWindow()
}

func (l *Line) constructEndLine() {
	dots := l.OriginLine.CopyDots()
	if l.IsReversed {
		for i := 0; i < len(dots); i++ {
			dots[i].X = -dots[i].X
		}
	}

}

func (l Line) isLineOutOfWindow() bool {
	var isOut = true

	for i := l.BeginLine.CurrentPosition; i < l.BeginLine.Size(); i++ {
		absDot, err := l.BeginLine.GetAbsoluteDot(i)
		if err != nil {
			log.Fatal(err)
		}
		isOut = absDot.Y < 0 || absDot.Y > window.height
		if !isOut {
			return isOut
		}
	}
	for i := 0; i < l.EndLine.CurrentPosition; i++ {
		absDot, err := l.BeginLine.GetAbsoluteDot(i)
		if err != nil {
			log.Fatal(err)
		}
		isOut = absDot.Y < 0 || absDot.Y > window.height
		if !isOut {
			return isOut
		}
	}

	return isOut
}
