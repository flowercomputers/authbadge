// Package authbadge draws the badge that pairs a CLI sign-in with the browser
// tab approving it: the user code on a plate, framed by a pattern seeded from
// the code. flop prints it in the terminal (ANSI) and spade serves it to the
// panel (SVG) from the same cells, so the two always match.
//
//	░▓▓░▓░▓▓░
//	▓░ABCDE░▓
//	░▓▓░▓░▓▓░
//
// It is deliberately unlike a vase avatar (17x9, yellow and purple): one green
// dithered like the terminal's shade blocks on the web, the terminal's own
// colors in flop, so a pending code is never mistaken for who you are.
package authbadge

import (
	"fmt"
	"hash/fnv"
	"strings"
)

const (
	Rows = 3
	// codeInset frame cells sit either side of the code on its row
	codeInset = 2
)

// Cols is the badge's width for a code of the given length.
func Cols(code string) int { return len(code) + 2*codeInset }

// Cell is one character cell of the badge: a code character, or a frame
// cell in one of two shades.
type Cell struct {
	Char  rune // 0 for a frame cell
	Shade int  // frame shade, 0 or 1
}

// Cells lays out the badge for code: Rows rows of Cols(code) cells. The frame
// is mirrored left to right and its shades come from the code's hash, so a
// different code looks different at a glance.
func Cells(code string) [][]Cell {
	h := fnv.New64a()
	_, _ = h.Write([]byte(code))
	bits := h.Sum64()

	cols := Cols(code)
	half := (cols + 1) / 2
	cells := make([][]Cell, Rows)
	for y := range cells {
		cells[y] = make([]Cell, cols)
		for x := 0; x < half; x++ {
			shade := int(bits & 1)
			bits >>= 1
			cells[y][x] = Cell{Shade: shade}
			cells[y][cols-1-x] = Cell{Shade: shade}
		}
	}
	for i := 0; i < len(code); i++ {
		cells[Rows/2][codeInset+i] = Cell{Char: rune(code[i])}
	}
	return cells
}

// ANSI renders the badge for a terminal, one line per row, in the terminal's
// own colors: the frame as block shades, the code bold on an inverted plate
// (or plain without color).
func ANSI(code string, color bool) string {
	shades := [2]rune{'▓', '░'}
	cells := Cells(code)
	lines := make([]string, len(cells))
	for y, row := range cells {
		var b strings.Builder
		for _, c := range row {
			switch {
			case c.Char == 0:
				b.WriteRune(shades[c.Shade])
			case color:
				fmt.Fprintf(&b, "\x1b[1;7m%c\x1b[0m", c.Char)
			default:
				b.WriteRune(c.Char)
			}
		}
		lines[y] = b.String()
	}
	return strings.Join(lines, "\n")
}

// one ink, the panel's accent, on paper, its card surface: the frame is
// dithered the way a terminal draws ▓ and ░, the plate solid ink like flop's
// inverted code. Both are classes so the badge follows the panel into dark
// mode (an <img> SVG sees the page's color scheme).
const (
	ink   = "ink"
	paper = "paper"
	style = `<style>.ink{fill:#418272}.paper{fill:#ffffff}` +
		`@media (prefers-color-scheme: dark){.ink{fill:#ffffff}.paper{fill:#222222}}</style>`
)

// shadeDots are the paper pixels of ▓ (and so the ink pixels of ░) in a 4x4
// tile: a staggered quarter, as terminal fonts draw the shade blocks
var shadeDots = [][2]int{{0, 0}, {2, 0}, {1, 2}, {3, 2}}

// SVG geometry: a terminal cell's proportions, scaled freely by the viewBox
const (
	cellW = 14
	cellH = 24
	font  = 18
)

// SVG renders the badge as a standalone image; same code, identical bytes.
func SVG(code string) string {
	cells := Cells(code)
	cols := Cols(code)
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" shape-rendering="crispEdges">`, cols*cellW, Rows*cellH)
	// the shades tile in badge space, so dots run on unbroken across cells
	b.WriteString(style)
	b.WriteString(`<defs>`)
	for shade, name := range shadeNames {
		dots, back := paper, ink // ▓: ink with paper dots
		if shade == 1 {
			dots, back = ink, paper // ░: paper with ink dots
		}
		fmt.Fprintf(&b, `<pattern id="%s" width="4" height="4" patternUnits="userSpaceOnUse"><rect width="4" height="4" class="%s"/>`, name, back)
		for _, d := range shadeDots {
			fmt.Fprintf(&b, `<rect x="%d" y="%d" width="1" height="1" class="%s"/>`, d[0], d[1], dots)
		}
		b.WriteString(`</pattern>`)
	}
	b.WriteString(`</defs>`)
	for y, row := range cells {
		for x, c := range row {
			fmt.Fprintf(&b, `<rect x="%d" y="%d" width="%d" height="%d" %s/>`, x*cellW, y*cellH, cellW, cellH, fill(c))
		}
	}
	fmt.Fprintf(&b, `<g font-family="ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', monospace" font-size="%d" font-weight="700" text-anchor="middle" dominant-baseline="central" text-rendering="geometricPrecision">`, font)
	for y, row := range cells {
		for x, c := range row {
			if c.Char != 0 {
				fmt.Fprintf(&b, `<text x="%d" y="%d" class="%s">%s</text>`, x*cellW+cellW/2, y*cellH+cellH/2, paper, escape(c.Char))
			}
		}
	}
	b.WriteString(`</g></svg>`)
	return b.String()
}

// shadeNames are the pattern ids for Shade 0 (▓) and 1 (░), as in ANSI
var shadeNames = [2]string{"dark", "light"}

// fill is a cell rect's paint attribute: the plate's ink, or a frame shade
func fill(c Cell) string {
	if c.Char != 0 {
		return `class="` + ink + `"`
	}
	return `fill="url(#` + shadeNames[c.Shade] + `)"`
}

func escape(r rune) string {
	switch r {
	case '<':
		return "&lt;"
	case '>':
		return "&gt;"
	case '&':
		return "&amp;"
	}
	return string(r)
}
