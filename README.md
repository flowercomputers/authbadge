# authbadge

Draws the badge that pairs a CLI sign-in with the browser tab approving it: the
user code on a plate, framed by a pattern seeded from the code.

```text
▓░▓▓▓▓▓░▓
░▓Q8D0N▓░
░▓▓░▓░▓▓░
```

The CLI prints it in the terminal and the server sends the same badge to the
browser as an SVG, so a person can check at a glance that the tab they are
approving belongs to the sign-in they started. [flop](https://flowercomputer.com),
Flower Computer's CLI, signs in this way.

```go
import "github.com/flowercomputers/authbadge"

fmt.Println(authbadge.ANSI("Q8D0N", true)) // the terminal: block shades, the code bold on an inverted plate
svg := authbadge.SVG("Q8D0N")              // the browser: a standalone image, following light and dark mode
cells := authbadge.Cells("Q8D0N")          // the layout itself, to draw it any other way
```

A code's badge never changes between versions: a CLI and a server built from
different versions still draw the same one (a test holds it to that).

No dependencies outside the standard library. Apache-2.0.
