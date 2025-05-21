package helpers

import (
	"path/filepath"

	xShared "xel/shared"

	"github.com/dev-kas/virtlang-go/v3/debugger"
)

func GenerateStackTrace(snapshot *debugger.Snapshot, cwd string) string {
	out := ""

	out += xShared.ColorPalette.Error.Sprintf("Stack trace (Most recent call first) (Stack Depth: %d)\n", len(snapshot.Stack))
	// The stack has most recent call last, so we reverse traverse it
	for i := len(snapshot.Stack) - 1; i >= 0; i-- {
		fname := snapshot.Stack[i].Filename
		name := snapshot.Stack[i].Name
		line := snapshot.Stack[i].Line

		relFname, err := filepath.Rel(cwd, fname)
		if err == nil {
			fname = relFname
		}

		prefix := "├─"
		if i == 0 {
			prefix = "└─"
		}

		out += xShared.ColorPalette.GrayMessage.Sprintf("%s at %s (%s:%d)\n", prefix, name, fname, line)
	}

	return out
}
