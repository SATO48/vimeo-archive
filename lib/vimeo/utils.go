package vimeo

import (
	"github.com/samber/lo"
	"github.com/silentsokolov/go-vimeo/vimeo"
)

func FindBestFile(files []*vimeo.File) *vimeo.File {
	original, found := lo.Find(files, func(d *vimeo.File) bool {
		return d.Quality == "source"
	})

	if !found {
		original = lo.MaxBy(files, func(a, b *vimeo.File) bool {
			return a.Width > b.Width
		})
	}

	return original
}
