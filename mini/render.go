package mini

import (
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/scraper"
	"github.com/vivi-app/vivi/style"
)

func renderMedia(media *scraper.Media) (rendered string) {
	rendered += media.String()

	if description := media.Description; description != "" {
		rendered += " " + style.Faint(description)
	}

	return
}

func renderExtension(ext *extension.Extension) (rendered string) {
	rendered += ext.String()

	if about := ext.Passport().About; about != "" {
		rendered += " " + style.Faint(about)
	}

	return
}
