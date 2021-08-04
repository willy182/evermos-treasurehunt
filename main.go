package main

import (
	"math/rand"
	"time"

	tl "github.com/JoelOtter/termloop"
)

const (
	maxX = 6
	maxY = 6

	minX = 0
	minY = 0
)

var (
	gridLayout [maxX][maxY]bool
	treasure   [maxX][maxY]bool
	clearPath  [maxX][maxY]*Path
)

type Player struct {
	*tl.Entity
}

// Tick
func (p *Player) Tick(ev tl.Event) {
	// Enable arrow key movement
	if ev.Type == tl.EventKey {
		x, y := p.Position()
		currX := x
		currY := y

		switch ev.Key {
		case tl.KeyArrowRight:
			x += 1
		case tl.KeyArrowLeft:
			x -= 1
		case tl.KeyArrowUp:
			y -= 1
		case tl.KeyArrowDown:
			y += 1
		case tl.KeyEnter:
			if treasure[x][y] {
				clearPath[x][y].SetText("$")
				clearPath[x][y].SetColor(tl.ColorGreen, tl.ColorDefault)
			}
		}

		// if obstacle then don't move
		if x < minX {
			x = minX
		}

		if y < minY {
			y = minY
		}

		if x > maxX {
			x = maxX
		}

		if y > maxY {
			y = maxY
		}

		if gridLayout[x][y] {
			p.SetPosition(currX, currY)
		} else {
			p.SetPosition(x, y)
		}

	}
}

// parsePlayer function for reading a Player out of JSON.
func parsePlayer(data map[string]interface{}) tl.Drawable {
	e := tl.NewEntity(
		data["x"].(int),
		data["y"].(int),
		1, 1,
	)
	e.SetCell(0, 0, &tl.Cell{
		Ch: []rune(data["ch"].(string))[0],
		Fg: tl.Attr(data["color"].(int)),
	})
	return &Player{e}
}

// Obstacle
type Obstacle struct {
	*tl.Text
}

func (o *Obstacle) Tick(ev tl.Event) {}

// Path
type Path struct {
	*tl.Text
}

func (p *Path) Tick(ev tl.Event) {}

// EventInfo
type EventInfo struct {
	*tl.Text
}

// NewEventInfo
func NewEventInfo(x, y int) *EventInfo {
	return &EventInfo{tl.NewText(x, y, "Yeay, Treasure!", tl.ColorWhite, tl.ColorBlack)}
}

// Info
type Info struct {
	*tl.Text
}

// NewInfo
func NewInfo(x, y int) *Info {
	return &Info{tl.NewText(x, y, "Press arrow key to move, Press enter to get the treasure", tl.ColorWhite, tl.ColorBlack)}
}

// isError
func isError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	gridLayout = [maxX][maxY]bool{
		{true, true, true, true, true, true},
		{true, false, false, false, false, true},
		{true, false, true, false, true, true},
		{true, false, false, false, false, true},
		{true, false, true, false, false, true},
		{true, true, true, true, true, true},
	}

	g := tl.NewGame()
	g.Screen().SetFps(30)

	l := tl.NewBaseLevel(tl.Cell{Bg: 0, Fg: 0})
	// add seed for random boolean treasure
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < maxX; i++ {
		for j := 0; j < maxY; j++ {
			if gridLayout[i][j] == true {
				// add obstacle
				g.Screen().AddEntity(&Obstacle{tl.NewText(j, i, "#", tl.ColorWhite, tl.ColorDefault)})
			} else {
				path := &Path{tl.NewText(j, i, ".", tl.ColorWhite, tl.ColorDefault)}
				// generate treasure
				treasure[j][i] = rand.Float32() > .5
				if treasure[j][i] {
					clearPath[j][i] = path

				}

				g.Screen().AddEntity(path)
			}
		}
	}

	// add player
	g.Screen().AddEntity(parsePlayer(map[string]interface{}{
		"x": 1, "y": 4, "ch": "x", "color": 100,
	}))
	// add info
	g.Screen().AddEntity(NewInfo(10, 1))

	g.Screen().SetLevel(l)
	g.Start()
}
