// Package font provides support for Allegro's font addon.
package font

// #include <allegro5/allegro.h>
// #include <allegro5/allegro_font.h>
// #include "../util.c"
import "C"
import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/phrasz/nag/allegro"
)

type Font C.ALLEGRO_FONT

type DrawFlags int

const (
	NO_KERNING              = C.ALLEGRO_NO_KERNING
	ALIGN_LEFT    DrawFlags = C.ALLEGRO_ALIGN_LEFT
	ALIGN_CENTRE            = C.ALLEGRO_ALIGN_CENTRE
	ALIGN_RIGHT             = C.ALLEGRO_ALIGN_RIGHT
	ALIGN_INTEGER           = C.ALLEGRO_ALIGN_INTEGER
)

// Initialise the font addon.
func Install() {
	C.al_init_font_addon()
}

// Shut down the font addon. This is done automatically at program exit, but
// can be called any time the user wishes as well.
func Uninstall() {
	C.al_shutdown_font_addon()
}

// Returns the (compiled) version of the addon, in the same format as
// al_get_allegro_version.
func Version() (major, minor, revision, release uint8) {
	v := uint32(C.al_get_allegro_font_version())
	major = uint8(v >> 24)
	minor = uint8((v >> 16) & 255)
	revision = uint8((v >> 8) & 255)
	release = uint8(v & 255)
	return
}

// Creates a monochrome bitmap font (8x8 pixels per character).
func Builtin() (*Font, error) {
	f := C.al_create_builtin_font()
	if f == nil {
		return nil, errors.New("failed to create builtin font")
	}
	return (*Font)(f), nil
}

// Loads a font from disk. This will use al_load_bitmap_font if you pass the
// name of a known bitmap format, or else al_load_ttf_font.
func LoadFont(filename string, size, flags int) (*Font, error) {
	filename_ := C.CString(filename)
	defer C.free_string(filename_)
	f := C.al_load_font(filename_, C.int(size), C.int(flags))
	if f == nil {
		return nil, fmt.Errorf("failed to load font '%s'", filename)
	}
	font := (*Font)(f)
	//runtime.SetFinalizer(font, font.Destroy)
	return font, nil
}

// Load a bitmap font from. It does this by first calling al_load_bitmap and
// then al_grab_font_from_bitmap. If you want to for example load an old A4
// font, you could load the bitmap yourself, then call al_convert_mask_to_alpha
// on it and only then pass it to al_grab_font_from_bitmap.
func LoadBitmapFont(filename string) (*Font, error) {
	filename_ := C.CString(filename)
	defer C.free_string(filename_)
	f := C.al_load_bitmap_font(filename_)
	if f == nil {
		return nil, fmt.Errorf("failed to load bitmap font '%s'", filename)
	}
	font := (*Font)(f)
	//runtime.SetFinalizer(font, font.Destroy)
	return font, nil
}

// Creates a new font from an Allegro bitmap. You can delete the bitmap after
// the function returns as the font will contain a copy for itself.
func GrabFontFromBitmap(bmp *allegro.Bitmap, ranges [][2]int) (*Font, error) {
	nRanges := len(ranges)
	n_ranges := nRanges * 2
	if n_ranges == 0 {
		return nil, errors.New("no ranges specified")
	}
	c_ranges := make([]C.int, n_ranges)
	//fmt.Printf("[nag][DEBUGGING] Current GFFBitmap range: 0 - %d\n", nRanges-1)
	for i := 0; i < len(ranges); i++ {
		for j := 0; j < 2; j++ { //len(ranges[i]); j++ {
			//fmt.Printf("[nag][DEBUGGING] Adding Element: - (%d,%d)\n", i, j)
			c_ranges[2*i+j] = C.int(ranges[i][j])
		}
	}

	f := C.al_grab_font_from_bitmap((*C.ALLEGRO_BITMAP)(unsafe.Pointer(bmp)), C.int(nRanges), (*C.int)(unsafe.Pointer(&c_ranges[0])))
	if f == nil {
		return nil, errors.New("failed to grab font from bitmap")
	}
	return (*Font)(f), nil
}

// Writes the NUL-terminated string text onto the target bitmap at position x,
// y, using the specified font.
func DrawText(font *Font, color allegro.Color, x, y float32, flags DrawFlags, text string) {
	text_ := C.CString(text)
	defer C.free_string(text_)
	C.al_draw_text((*C.ALLEGRO_FONT)(font),
		*((*C.ALLEGRO_COLOR)(unsafe.Pointer(&color))), // is there an easier way to get this converted?
		C.float(x),
		C.float(y),
		C.int(flags),
		text_)
}

// Like al_draw_text, but justifies the string to the region x1-x2.
func DrawJustifiedText(font *Font, color allegro.Color, x1, x2, y, diff float32, flags DrawFlags, text string) {
	text_ := C.CString(text)
	defer C.free_string(text_)
	C.al_draw_justified_text((*C.ALLEGRO_FONT)(font),
		*((*C.ALLEGRO_COLOR)(unsafe.Pointer(&color))), // is there an easier way to get this converted?
		C.float(x1),
		C.float(x2),
		C.float(y),
		C.float(diff),
		C.int(flags),
		text_)
}

func DrawTextf(font *Font, color allegro.Color, x, y float32, flags DrawFlags, format string, a ...interface{}) {
	// C.al_draw_textf
	text_ := C.CString(fmt.Sprintf(format, a...))
	defer C.free_string(text_)
	C.al_draw_text((*C.ALLEGRO_FONT)(font),
		*((*C.ALLEGRO_COLOR)(unsafe.Pointer(&color))), // is there an easier way to get this converted?
		C.float(x),
		C.float(y),
		C.int(flags),
		text_)
}

func DrawJustifiedTextf(font *Font, color allegro.Color, x1, x2, y, diff float32, flags DrawFlags, format string, a ...interface{}) {
	// C.al_draw_justified_textf
	text_ := C.CString(fmt.Sprintf(format, a))
	defer C.free_string(text_)
	C.al_draw_justified_text((*C.ALLEGRO_FONT)(font),
		*((*C.ALLEGRO_COLOR)(unsafe.Pointer(&color))), // is there an easier way to get this converted?
		C.float(x1),
		C.float(x2),
		C.float(y),
		C.float(diff),
		C.int(flags),
		text_)
}

// Frees the memory being used by a font structure. Does nothing if passed NULL.
func (f *Font) Destroy() {
	C.al_destroy_font((*C.ALLEGRO_FONT)(f))
}

// Returns the usual height of a line of text in the specified font. For bitmap
// fonts this is simply the height of all glyph bitmaps. For truetype fonts it
// is whatever the font file specifies. In particular, some special glyphs may
// be higher than the height returned here.
func (f *Font) LineHeight() int {
	return int(C.al_get_font_line_height((*C.ALLEGRO_FONT)(f)))
}

// Returns the ascent of the specified font.
func (f *Font) Ascent() int {
	return int(C.al_get_font_ascent((*C.ALLEGRO_FONT)(f)))
}

// Returns the descent of the specified font.
func (f *Font) Descent() int {
	return int(C.al_get_font_descent((*C.ALLEGRO_FONT)(f)))
}

// Calculates the length of a string in a particular font, in pixels.
func (f *Font) TextWidth(text string) int {
	text_ := C.CString(text)
	defer C.free_string(text_)
	return int(C.al_get_text_width((*C.ALLEGRO_FONT)(f), text_))
}

// Sometimes, the al_get_text_width and al_get_font_line_height functions are
// not enough for exact text placement, so this function returns some
// additional information.
func (f *Font) TextDimensions(text string) (bbx, bby, bbw, bbh int) {
	var cbbx, cbby, cbbw, cbbh C.int
	text_ := C.CString(text)
	defer C.free_string(text_)
	C.al_get_text_dimensions((*C.ALLEGRO_FONT)(f), text_,
		&cbbx, &cbby, &cbbw, &cbbh)
	return int(cbbx), int(cbby), int(cbbw), int(cbbh)
}

//From Font:
//static int color_get_glyph_advance(ALLEGRO_FONT const *f, int codepoint1, int codepoint2)
//static bool color_get_glyph_dimensions(ALLEGRO_FONT const *f, int codepoint, int *bbx, int *bby, int *bbw, int *bbh)

//From Text.c
//int al_get_glyph_width(const ALLEGRO_FONT *f, int codepoint)
//bool al_get_glyph_dimensions(const ALLEGRO_FONT *f, int codepoint, int *bbx, int *bby, int *bbw, int *bbh)
//bool al_get_glyph(const ALLEGRO_FONT *f, int prev_codepoint, int codepoint, ALLEGRO_GLYPH *glyph)
//
//

//GetGlyphAdvance returns the width of next glyph?
//int al_get_glyph_advance(const ALLEGRO_FONT *f, int codepoint1, int codepoint2)
//int al_get_glyph_advance(const ALLEGRO_FONT *f, int codepoint1, int codepoint2)
func GetGlyphAdvance(font *Font, codepoint1 int, codepoint2 int) int {
	return int(C.al_get_glyph_advance((*C.ALLEGRO_FONT)(font), C.int(codepoint1), C.int(codepoint2)))
}

//DrawGlyph draws a glyph
//int al_get_glyph_advance(const ALLEGRO_FONT *f, int codepoint1, int codepoint2)
//void al_draw_glyph(const ALLEGRO_FONT *f, ALLEGRO_COLOR color, float x, float y, int codepoint)
func DrawGlyph(font *Font, color allegro.Color, x float32, y float32, codepoint int) {
	C.al_draw_glyph((*C.ALLEGRO_FONT)(font), *((*C.ALLEGRO_COLOR)(unsafe.Pointer(&color))), C.float(x), C.float(y), C.int(codepoint))
	return
}

//font *Font, color allegro.Color, x1, x2, y, diff float32, flags DrawFlags, format string, a ...interface{}
