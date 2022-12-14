package passport

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/awirix/awirix/color"
	"github.com/awirix/awirix/executil"
	"github.com/awirix/awirix/icon"
	"github.com/awirix/awirix/style"
	"runtime"
	"strings"
)

type Requirements struct {
	OS       map[string]bool `json:"os,omitempty"`
	Programs []string        `json:"programs,omitempty"`
}

func (d *Requirements) Info() string {
	var b strings.Builder

	if len(d.OS) != 0 {
		b.WriteString("OS ")

		if enabled, ok := d.OS[runtime.GOOS]; ok && enabled {
			b.WriteString(style.Fg(color.Green)(icon.Check))
		} else {
			b.WriteString(style.Fg(color.Red)(icon.Cross))
		}

		b.WriteString(style.Faint(" Available on: "))
		b.WriteString(style.Faint(strings.Join(lo.Keys(d.OS), ", ")))
		b.WriteString(style.Faint(fmt.Sprintf(". You're on: %s", runtime.GOOS)))

		b.WriteRune('\n')
	}

	if len(d.Programs) != 0 {
		b.WriteString("Programs ")

		for _, program := range d.Programs {
			if executil.ProgramInPath(program) {
				b.WriteString(style.Fg(color.Green)(icon.Check))
			} else {
				b.WriteString(style.Fg(color.Red)(icon.Cross))
			}

			b.WriteRune(' ')
			b.WriteString(style.Faint(program))
			b.WriteRune(' ')
		}
	}

	return b.String()
}

func (d *Requirements) Matches() bool {
	if enabled, ok := d.OS[runtime.GOOS]; len(d.OS) != 0 && ok && !enabled {
		return false
	}

	for _, program := range d.Programs {
		if !executil.ProgramInPath(program) {
			return false
		}
	}

	return true
}
