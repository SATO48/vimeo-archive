package vimeo

import (
	"context"
	"errors"
	"log/slog"

	"github.com/silentsokolov/go-vimeo/vimeo"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type Scraper struct {
	api         *vimeo.Client
	pageSize    int
	maxPages    int
	pagePointer int
	totalPages  int
}

type OptionFunc func(*Scraper)

func WithPageSize(size int) OptionFunc {
	return func(vs *Scraper) {
		vs.pageSize = size
	}
}

func WithMaxPages(pages int) OptionFunc {
	return func(vs *Scraper) {
		vs.maxPages = pages
	}
}

func WithPagePointer(pointer int) OptionFunc {
	return func(vs *Scraper) {
		vs.pagePointer = pointer
	}
}

func NewScraper(options ...OptionFunc) *Scraper {
	vs := &Scraper{}
	vs.pageSize = 25

	for _, opt := range options {
		opt(vs)
	}

	vs.api = vimeo.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: viper.GetString("VIMEO_ACCESS_KEY")},
	)), nil)

	return vs
}

func (vs *Scraper) HasNextPage() bool {
	if vs.pagePointer == 0 {
		return true
	}

	if vs.pagePointer >= vs.maxPages {
		return false
	}

	if vs.pagePointer >= vs.totalPages {
		return false
	}

	return true
}

func (vs *Scraper) ListVideos() ([]*vimeo.Video, error) {
	vs.pagePointer++
	slog.Info("listing videos from vimeo", "page", vs.pagePointer, "total", vs.totalPages)

	videos, r, err := vs.api.Users.ListVideo("", vimeo.OptPerPage(vs.pageSize), vimeo.OptPage(vs.pagePointer))
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, errors.New("response is nil")
	}

	vs.totalPages = r.TotalPages

	return videos, nil
}

func (vs *Scraper) GetVideo(id int) (*vimeo.Video, error) {
	video, _, err := vs.api.Videos.Get(id)
	if err != nil {
		return nil, err
	}

	return video, nil
}
