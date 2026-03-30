package console

import (
	"fmt"
	"io"
)

type UI struct {
	out io.Writer
	err io.Writer
}

func New(out, err io.Writer) *UI {
	return &UI{
		out: out,
		err: err,
	}
}

func (ui *UI) PrintStatus(format string, args ...any) {
	_, _ = fmt.Fprintf(ui.out, format+"\n", args...)
}

func (ui *UI) PrintError(format string, args ...any) {
	_, _ = fmt.Fprintf(ui.err, format+"\n", args...)
}
