package vimeo

import (
	"vimeo-archive/lib/model"

	"github.com/davecgh/go-spew/spew"
	"github.com/silentsokolov/go-vimeo/v2/vimeo"
)

type Archiver struct {
	vc  *vimeo.Client
	vb  *model.VideoBox
	max uint64
}

type ArchiverOptionFunc func(*Archiver)

func WithVimeoClient(vc *vimeo.Client) ArchiverOptionFunc {
	return func(a *Archiver) {
		a.vc = vc
	}
}

func WithVideoBox(vb *model.VideoBox) ArchiverOptionFunc {
	return func(a *Archiver) {
		a.vb = vb
	}
}

func WithMax(max uint64) ArchiverOptionFunc {
	return func(a *Archiver) {
		a.max = max
	}
}

func NewArchiver(options ...ArchiverOptionFunc) *Archiver {
	a := &Archiver{}

	for _, opt := range options {
		opt(a)
	}

	return a
}

func (a *Archiver) Archive() error {
	videos, err := a.vb.Query(
		model.Video_.DownloadedTime.IsNil(),
	).Limit(a.max).Find()
	if err != nil {
		return err
	}

	for _, video := range videos {
		if err := a.archiveVideo(video); err != nil {
			return err
		}
	}

	return nil
}

func (a *Archiver) archiveVideo(v *model.Video) error {
	r, _, err := a.vc.Videos.Get(int(v.Id))
	if err != nil {
		return err
	}

	f := FindBestFile(r.Files)

	spew.Dump(f.Link)

	return nil
}
