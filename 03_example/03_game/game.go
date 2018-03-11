// This example demonstrates a very crude tile-based game.
package main

import (
	"fmt"
	"math"
	"os"

	"github.com/phrasz/nag/allegro"
)

const titleSize = 30
const gopherSize = 20
const startX = 6
const startY = 6
const gopherSpeed = 6
const fps = 30

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
	displayWidth  = titleSize * len(gameMap[0])
	displayHeight = titleSize * len(gameMap)
)

//Game Global game object.
type Game struct {
	gopher     *Gopher
	background *allegro.Bitmap
	keyboard   *allegro.KeyboardState
	tiles      map[int]*Tile
}

//RenderTile renders the tile with the given id at the given position.
func (game *Game) RenderTile(tile, x, y int) {
	t, ok := game.tiles[tile]
	if !ok {
		return
	}
	t.Render(x, y)
}

//Entity Something visible.
type Entity struct {
	image *allegro.Bitmap
}

//Object An object with x- and y- coordinates along with a width
// and a height.
type Object struct {
	Entity
	x, y, w, h float32
}

//Render draws the object on to the target bitmap.
func (ob *Object) Render() {
	ob.image.Draw(ob.x, ob.y, allegro.FLIP_NONE)
}

//Move moves the object by (x,y), but not letting it move through
// any tiles except those with id 0. This is currently very dumb, but
// gets the job done for the purposes of this example.
func (ob *Object) Move(x, y float32) {
	var (
		tx = ob.x + x
		ty = ob.y + y
	)
	if x > 0 {
		xtile := int(math.Floor(float64((tx + ob.w) / titleSize)))
		ytile1 := int(math.Floor(float64(ty / titleSize)))
		ytile2 := int(math.Floor(float64((ty + ob.h) / titleSize)))
		if gameMap[ytile1][xtile] == 0 && gameMap[ytile2][xtile] == 0 {
			ob.x = tx
		}
	} else if x < 0 {
		xtile := int(math.Floor(float64(tx / titleSize)))
		ytile1 := int(math.Floor(float64(ty / titleSize)))
		ytile2 := int(math.Floor(float64((ty + ob.h) / titleSize)))
		if gameMap[ytile1][xtile] == 0 && gameMap[ytile2][xtile] == 0 {
			ob.x = tx
		}
	} else if y > 0 {
		xtile1 := int(math.Floor(float64(tx / titleSize)))
		xtile2 := int(math.Floor(float64((tx + ob.w) / titleSize)))
		ytile := int(math.Floor(float64((ty + ob.h) / titleSize)))
		if gameMap[ytile][xtile1] == 0 && gameMap[ytile][xtile2] == 0 {
			ob.y = ty
		}
	} else if y < 0 {
		xtile1 := int(math.Floor(float64(tx / titleSize)))
		xtile2 := int(math.Floor(float64((tx + ob.w) / titleSize)))
		ytile := int(math.Floor(float64(ty / titleSize)))
		if gameMap[ytile][xtile1] == 0 && gameMap[ytile][xtile2] == 0 {
			ob.y = ty
		}
	}
}

//Tile A tile on the field.
type Tile struct {
	Entity
	id int
}

//Render renders this tile at the given position.
func (t *Tile) Render(x, y int) {
	t.image.Draw(float32(x*titleSize), float32(y*titleSize), allegro.FLIP_NONE)
}

//Gopher The character.
type Gopher struct {
	Object
}

//Update is called once every frame, and should take care of handling
// updates to the game world.
func (game *Game) Update() {
	game.keyboard.Get()
	if game.keyboard.IsDown(allegro.KEY_RIGHT) {
		game.gopher.Move(gopherSpeed, 0)
	} else if game.keyboard.IsDown(allegro.KEY_LEFT) {
		game.gopher.Move(-gopherSpeed, 0)
	} else if game.keyboard.IsDown(allegro.KEY_DOWN) {
		game.gopher.Move(0, gopherSpeed)
	} else if game.keyboard.IsDown(allegro.KEY_UP) {
		game.gopher.Move(0, -gopherSpeed)
	}
}

//Render draws everything to the screen.
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
			running    = true
			err        error
		)

		game := new(Game)
		game.tiles = make(map[int]*Tile)
		game.keyboard = new(allegro.KeyboardState)

		if eventQueue, err = allegro.CreateEventQueue(); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			return
		}
		defer eventQueue.Destroy()

		if err := allegro.InstallKeyboard(); err != nil {
			panic(err)
		}

		allegro.SetNewDisplayFlags(allegro.WINDOWED)
		if display, err = allegro.CreateDisplay(displayWidth, displayHeight); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			return
		}
		defer display.Destroy()
		display.SetWindowTitle("You Can't Leave!")

		screen := allegro.TargetBitmap()

		game.gopher = new(Gopher)
		game.gopher.image = allegro.CreateBitmap(gopherSize, gopherSize)
		allegro.SetTargetBitmap(game.gopher.image)
		allegro.ClearToColor(allegro.MapRGB(0xFF, 0, 0))

		game.gopher.w = float32(game.gopher.image.Width())
		game.gopher.h = float32(game.gopher.image.Height())
		game.gopher.x = float32((startX * titleSize) - (game.gopher.w / 2))
		game.gopher.y = float32((startY * titleSize) - (game.gopher.h / 2))

		whiteTile := &Tile{id: 0}
		whiteTile.image = allegro.CreateBitmap(titleSize, titleSize)
		allegro.SetTargetBitmap(whiteTile.image)
		allegro.ClearToColor(allegro.MapRGB(0xFF, 0xFF, 0xFF))
		game.tiles[0] = whiteTile

		blackTile := &Tile{id: 1}
		blackTile.image = allegro.CreateBitmap(titleSize, titleSize)
		allegro.SetTargetBitmap(blackTile.image)
		allegro.ClearToColor(allegro.MapRGB(0, 0, 0))
		game.tiles[1] = blackTile

		// create the background
		game.background = allegro.CreateBitmap(displayWidth, displayHeight)
		allegro.SetTargetBitmap(game.background)
		allegro.HoldBitmapDrawing(true)
		for y, row := range gameMap {
			for x, tile := range row {
				game.RenderTile(tile, x, y)
			}
		}
		allegro.HoldBitmapDrawing(false)
		allegro.SetTargetBitmap(screen)

		timer, err := allegro.CreateTimer(1.0 / fps)
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
