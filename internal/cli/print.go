package cli

import (
	"fmt"
	"io"
)

// Writes to a command's output stream never meaningfully fail (the stream is a
// terminal or buffer), so these helpers centralise discarding that error and
// keep the call sites readable.

func fprintf(w io.Writer, format string, a ...any) {
	_, _ = fmt.Fprintf(w, format, a...)
}

func fprintln(w io.Writer, a ...any) {
	_, _ = fmt.Fprintln(w, a...)
}

func fprint(w io.Writer, a ...any) {
	_, _ = fmt.Fprint(w, a...)
}
