package authbadge

import (
	"strings"
	"testing"
)

func TestCellsLayout(t *testing.T) {
	cells := Cells("Q8D0N")
	if len(cells) != Rows {
		t.Fatalf("rows = %d, want %d", len(cells), Rows)
	}
	for y, row := range cells {
		if len(row) != 9 {
			t.Fatalf("row %d has %d cells, want 9", y, len(row))
		}
		// the frame is mirrored
		for x := range row {
			if row[x].Char == 0 && row[x] != row[len(row)-1-x] {
				t.Errorf("row %d not mirrored at %d", y, x)
			}
		}
	}
	if got := string(runes(cells[1][2:7])); got != "Q8D0N" {
		t.Errorf("middle row code = %q", got)
	}
}

func TestDeterministicAndDistinct(t *testing.T) {
	if SVG("ABCDE") != SVG("ABCDE") {
		t.Fatal("same code, different SVG")
	}
	if ANSI("ABCDE", true) == ANSI("ABCDF", true) {
		t.Fatal("different codes, same badge")
	}
}

func TestPlainANSI(t *testing.T) {
	lines := strings.Split(ANSI("ABCDE", false), "\n")
	if len(lines) != Rows || !strings.Contains(lines[1], "ABCDE") {
		t.Fatalf("plain badge = %q", lines)
	}
	for _, l := range lines {
		if strings.Contains(l, "\x1b") {
			t.Fatalf("plain badge has escapes: %q", l)
		}
	}
}

func runes(cells []Cell) []rune {
	out := make([]rune, len(cells))
	for i, c := range cells {
		out[i] = c.Char
	}
	return out
}

// a code's badge never changes: a CLI and a server built from different
// versions of this package must still draw the same one
func TestACodesBadgeNeverChanges(t *testing.T) {
	for code, want := range map[string]string{
		"Q8D0N": "▓░▓▓▓▓▓░▓\n░▓Q8D0N▓░\n░▓▓░▓░▓▓░",
		"ABCDE": "▓▓▓░▓░▓▓▓\n▓░ABCDE░▓\n░▓▓░░░▓▓░",
	} {
		if got := ANSI(code, false); got != want {
			t.Errorf("%s:\n%s\nwant:\n%s", code, got, want)
		}
	}
}
