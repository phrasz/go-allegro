// This example demonstrates a very crude tile-based game.
package main

import (
	"fmt"
	"math"
	"os"

	"github.com/phrasz/nag/allegro"
)

const TILE_SIZE = 30
const GOPHER_SIZE = 20
const START_X = 6
const START_Y = 6
const GOPHER_SPEED = 6
const FPS = 30

var gameMap = [][]int{
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
}

var (
	DISPLAY_WIDTH  = TILE_SIZE * len(gameMap[0])
	DISPLAY_HEIGHT = TILE_SIZE * len(gameMap)
)

// Global game object.
type Game struct {
	gopher     *Gopher
	background *allegro.Bitmap
	keyboard   *allegro.KeyboardState
	tiles      map[int]*Tile
}

// RenderTile() renders the tile with the given id at the given position.
func (g *Game) RenderTile(tile, x, y int) {
	t, ok := g.tiles[tile]
	if !ok {
		return
	}
	t.Render(x, y)
}

// Something visible.
type Entity struct {
	image *allegro.Bitmap
}

// An object with x- and y- coordinates along with a width
// and a height.
type Object struct {
	Entity
	x, y, w, h float32
}

// Render() draws the object on to the target bitmap.
func (ob *Object) Render() {
	ob.image.Draw(ob.x, ob.y, allegro.FLIP_NONE)
}

// Move() moves the object by (x,y), but not letting it move through
// any tiles except those with id 0. This is currently very dumb, but
// gets the job done for the purposes of this example.
func (ob *Object) Move(x, y float32) {
	var (
		tx float32 = ob.x + x
		ty float32 = ob.y + y
	)
	if x > 0 {
		xtile := int(math.Floor(float64((tx + ob.w) / TILE_SIZE)))
		ytile1 := int(math.Floor(float64(ty / TILE_SIZE)))
		ytile2 := int(math.Floor(float64((ty + ob.h) / TILE_SIZE)))
		if gameMap[ytile1][xtile] == 0 && gameMap[ytile2][xtile] == 0 {
			ob.x = tx
		}
	} else if x < 0 {
		xtile := int(math.Floor(float64(tx / TILE_SIZE)))
		ytile1 := int(math.Floor(float64(ty / TILE_SIZE)))
		ytile2 := int(math.Floor(float64((ty + ob.h) / TILE_SIZE)))
		if gameMap[ytile1][xtile] == 0 && gameMap[ytile2][xtile] == 0 {
			ob.x = tx
		}
	} else if y > 0 {
		xtile1 := int(math.Floor(float64(tx / TILE_SIZE)))
		xtile2 := int(math.Floor(float64((tx + ob.w) / TILE_SIZE)))
		ytile := int(math.Floor(float64((ty + ob.h) / TILE_SIZE)))
		if gameMap[ytile][xtile1] == 0 && gameMap[ytile][xtile2] == 0 {
			ob.y = ty
		}
	} else if y < 0 {
		xtile1 := int(math.Floor(float64(tx / TILE_SIZE)))
		xtile2 := int(math.Floor(float64((tx + ob.w) / TILE_SIZE)))
		ytile := int(math.Floor(float64(ty / TILE_SIZE)))
		if gameMap[ytile][xtile1] == 0 && gameMap[ytile][xtile2] == 0 {
			ob.y = ty
		}
	}
}

// A tile on the field.
type Tile struct {
	Entity
	id int
}

// Render() renders this tile at the given position.
func (t *Tile) Render(x, y int) {
	t.image.Draw(float32(x*TILE_SIZE), float32(y*TILE_SIZE), allegro.FLIP_NONE)
}

// The character.
type Gopher struct {
	Object
}

// Update() is called once every frame, and should take care of handling
// updates to the game world.
func (game *Game) Update() {
	game.keyboard.Get()
	if game.keyboard.IsDown(allegro.KEY_RIGHT) {
		game.gopher.Move(GOPHER_SPEED, 0)
	} else if game.keyboard.IsDown(allegro.KEY_LEFT) {
		game.gopher.Move(-GOPHER_SPEED, 0)
	} else if game.keyboard.IsDown(allegro.KEY_DOWN) {
		game.gopher.Move(0, GOPHER_SPEED)
	} else if game.keyboard.IsDown(allegro.KEY_UP) {
		game.gopher.Move(0, -GOPHER_SPEED)
	}
}

// Render() draws everything to the screen.
func (game *Game) Render() {
	allegro.ClearToColor(allegro.MapRGB(0, 0, 0))
	allegro.HoldBitmapDrawing(true)
	game.background.Draw(0, 0, allegro.FLIP_NONE)
	game.gopher.Render()
	allegro.HoldBitmapDrawing(false)
	allegro.FlipDisplay()
}

func main() {
	allegro.Run(func() {
		var (
			display    *allegro.Display
			eventQueue *allegro.EventQueue
			running    bool = true
			err        error
		)

		game := new(Game)
		game.tiles = make(map[int]*Tile)
		game.keyboard = new(allegro.KeyboardState)

		if eventQueue, err = allegro.CreateEventQueue(); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			return
		} else {
			defer eventQueue.Destroy()
		}

		if err := allegro.InstallKeyboard(); err != nil {
			panic(err)
		}

		allegro.SetNewDisplayFlags(allegro.WINDOWED)
		if display, err = allegro.CreateDisplay(DISPLAY_WIDTH, DISPLAY_HEIGHT); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			return
		} else {
			defer display.Destroy()
			display.SetWindowTitle("You Can't Leave!")
		}

		screen := allegro.TargetBitmap()

		game.gopher = new(Gopher)
		game.gopher.image = allegro.CreateBitmap(GOPHER_SIZE, GOPHER_SIZE)
		allegro.SetTargetBitmap(game.gopher.image)
		allegro.ClearToColor(allegro.MapRGB(0xFF, 0, 0))

		game.gopher.w = float32(game.gopher.image.Width())
		game.gopher.h = float32(game.gopher.image.Height())
		game.gopher.x = float32((START_X * TILE_SIZE) - (game.gopher.w / 2))
		game.gopher.y = float32((START_Y * TILE_SIZE) - (game.gopher.h / 2))

		whiteTile := &Tile{id: 0}
		whiteTile.image = allegro.CreateBitmap(TILE_SIZE, TILE_SIZE)
		allegro.SetTargetBitmap(whiteTile.image)
		allegro.ClearToColor(allegro.MapRGB(0xFF, 0xFF, 0xFF))
		game.tiles[0] = whiteTile

		blackTile := &Tile{id: 1}
		blackTile.image = allegro.CreateBitmap(TILE_SIZE, TILE_SIZE)
		allegro.SetTargetBitmap(blackTile.image)
		allegro.ClearToColor(allegro.MapRGB(0, 0, 0))
		game.tiles[1] = blackTile

		// create the background
		game.background = allegro.CreateBitmap(DISPLAY_WIDTH, DISPLAY_HEIGHT)
		allegro.SetTargetBitmap(game.background)
		allegro.HoldBitmapDrawing(true)
		for y, row := range gameMap {
			for x, tile := range row {
				game.RenderTile(tile, x, y)
			}
		}
		allegro.HoldBitmapDrawing(false)
		allegro.SetTargetBitmap(screen)

		timer, err := allegro.CreateTimer(1.0 / FPS)
		if err != nil {
			panic(err)
		}

		eventQueue.Register(display)
		eventQueue.Register(timer)

		redraw := false
		timer.Start()

		var event allegro.Event
		for running {
			eventQueue.WaitForEvent(&event)
			switch eventQueue.WaitForEvent(&event).(type) {
			case allegro.TimerEvent:
				redraw = true

			case allegro.DisplayCloseEvent:
				println("display close")
				running = false

			default:
				// unknown event
			}

			if !running {
				break
			}

			if redraw && eventQueue.IsEmpty() {
				game.Update()
				game.Render()
				redraw = false
			}
		}
	})
}
