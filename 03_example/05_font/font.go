package main

import (
	"fmt"
	"math"

	"github.com/phrasz/nag/allegro"
	"github.com/phrasz/nag/allegro/font"
	"github.com/phrasz/nag/allegro/image"
)

const fps int = 60
const euro = "€"

func waitForEsc(display *allegro.Display) {
	var (
		keyboard    *allegro.EventSource
		eventQueue  *allegro.EventQueue
		screenClone *allegro.Bitmap
		err         error
	)

	allegro.InstallKeyboard()

	if eventQueue, err = allegro.CreateEventQueue(); err == nil {
		defer eventQueue.Destroy()
	} else {
		panic(err)
	}
	keyboard, _ = allegro.KeyboardEventSource()
	eventQueue.RegisterEventSource(keyboard)

	eventQueue.Register(display)

	screenClone, _ = (allegro.TargetBitmap()).Clone()
	for {
		var eventWait allegro.Event

		event := eventQueue.WaitForEvent(&eventWait)

		switch event.(type) {
		case allegro.DisplayCloseEvent:
			return
		case allegro.KeyDownEvent:
			if event.(allegro.KeyDownEvent).KeyCode() == allegro.KEY_ESCAPE {
				return
			}
		case allegro.DisplayExposeEvent:
			x := event.(allegro.DisplayExposeEvent).X()
			y := event.(allegro.DisplayExposeEvent).Y()
			w := event.(allegro.DisplayExposeEvent).Width()
			h := event.(allegro.DisplayExposeEvent).Height()

			screenClone.DrawRegion(float32(x), float32(y), float32(w), float32(h), float32(x), float32(y), 0)

			allegro.UpdateDisplayRegion(x, y, w, h)
		}
	}
	screenClone.Destroy()
	eventQueue.Destroy()
}

func main() {
	allegro.Run(func() {
		var (
			display    *allegro.Display
			bitmap     *allegro.Bitmap
			fontBitmap *allegro.Bitmap
			f1         *font.Font
			f2         *font.Font
			f3         *font.Font

			rangeVal int
			index    int
			x        int
			y        int

			err error
		)

		ranges := [][2]int{
			[2]int{0x0020, 0x007F}, // ASCII
			[2]int{0x00A1, 0x00FF}, // Latin 1
			[2]int{0x0100, 0x017F}, // Extended-A
			[2]int{0x20AC, 0x20AC}, // Euro
		}

		if err = image.Install(); err != nil {
			panic(err)
		}

		font.Install()

		// Note: C example's common.c; Android Touch is installed
		//    al_install_touch_input();
		//    al_android_set_apk_file_interface();

		allegro.SetNewDisplayOption(allegro.SINGLE_BUFFER, 1, allegro.SUGGEST)

		allegro.SetNewDisplayFlags(allegro.GENERATE_EXPOSE_EVENTS)

		display, err = allegro.CreateDisplay(640, 480)
		if err != nil {
			fmt.Printf("\n[ERROR] Cannot create a display!")
			return
		}

		bitmap, err = allegro.LoadBitmap("data/mysha.pcx")
		if err != nil {
			fmt.Printf("\n[ERROR] Failed to load mysha.pcx!\n")
		}

		f1, err = font.LoadFont("data/bmpfont.tga", 0, 0)
		if err != nil {
			fmt.Printf("\n[ERROR] Failed to load bmpfont.tga!\n")
		}

		fontBitmap, err = allegro.LoadBitmap("data/a4_font.tga")
		if err != nil {
			fmt.Printf("\n[ERROR] Failed to load a4_font.tga\n")
		}

		f2, err = font.GrabFontFromBitmap(fontBitmap, ranges)
		fmt.Printf("[DEBUGGING] Ranges vals: %d\n", ranges)
		f3, err = font.Builtin()
		if err != nil {
			fmt.Printf("\n[ERROR] Failed to create builtin font\n")
		}

		bitmap.DrawScaled(0, 0, 320, 240, 0, 0, 640, 480, 0)

		font.DrawTextf(f1, allegro.MapRGB(255, 0, 0), 10, 10, 0, "red")
		font.DrawTextf(f1, allegro.MapRGB(0, 255, 0), 120, 10, 0, "green")

		font.DrawTextf(f2, allegro.MapRGB(0, 0, 255), 60, 60, 0, "Mysha's 0.02"+euro)

		font.DrawTextf(f3, allegro.MapRGB(255, 255, 0), 20, 200, font.ALIGN_CENTRE, "a string from builtin font data")

		x = 10
		y = 300

		// Draw all individual glyphs the f2 font's range in rainbow colors.
		font.DrawTextf(f3, allegro.MapRGB(0, 255, 255), float32(x), float32(y-20), 0, "Draw glyphs: ")
		for rangeVal = 0; rangeVal < len(ranges); rangeVal++ {
			start := int(ranges[rangeVal][0])
			stop := ranges[rangeVal][1]

			for index = start; index < stop; index++ {
				width := font.GetGlyphAdvance(f2, index, font.NO_KERNING)

				r := byte(math.Abs(math.Sin(math.Pi*float64((index)*36)/360.0)) * 255.0)
				g := byte(math.Abs(math.Sin(math.Pi*float64((index+12)*36)/360.0)) * 255.0)
				b := byte(math.Abs(math.Sin(math.Pi*float64((index+24)*36)/360.0)) * 255.0)

				font.DrawGlyph(f2, allegro.MapRGB(r, g, b), float32(x), float32(y), index)
				x += width
				if x > (display.Width() - 10) {
					x = 10
					y += f2.LineHeight()
				}
			}
		}

		allegro.FlipDisplay()

		waitForEsc(display)
		bitmap.Destroy()
		fontBitmap.Destroy()

		f1.Destroy()
		f2.Destroy()
		f3.Destroy()
		return
	})
}
