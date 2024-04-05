package vimeo

import (
	"errors"
	"log/slog"

	"github.com/silentsokolov/go-vimeo/v2/vimeo"
)

type Scraper struct {
	api         *vimeo.Client
	pageSize    int
	maxPages    int
	pagePointer int
	total       int
}

type ScraperOptionFunc func(*Scraper)

func WithAPI(api *vimeo.Client) ScraperOptionFunc {
	return func(vs *Scraper) {
		vs.api = api
	}
}

func WithPageSize(size int) ScraperOptionFunc {
	return func(vs *Scraper) {
		vs.pageSize = size
	}
}

func WithMaxPages(pages int) ScraperOptionFunc {
	return func(vs *Scraper) {
		vs.maxPages = pages
	}
}

func WithPagePointer(pointer int) ScraperOptionFunc {
	return func(vs *Scraper) {
		vs.pagePointer = pointer
	}
}

func NewScraper(options ...ScraperOptionFunc) *Scraper {
	vs := &Scraper{}
	vs.pageSize = 25

	for _, opt := range options {
		opt(vs)
	}

	return vs
}

func (vs *Scraper) HasNextPage() bool {
	if vs.total == 0 {
		return true
	}

	if vs.maxPages > 0 && vs.pagePointer >= vs.maxPages {
		return false
	}

	if vs.pagePointer >= vs.total/vs.pageSize+1 {
		return false
	}

	return true
}

func (vs *Scraper) ListVideos() ([]*vimeo.Video, error) {
	vs.pagePointer++
	slog.Info("listing videos from vimeo", "page", vs.pagePointer, "total", vs.total)

	videos, r, err := vs.api.Users.ListVideo("",
		vimeo.OptPerPage(vs.pageSize),
		vimeo.OptPage(vs.pagePointer),
		vimeo.OptSort("date"),
		vimeo.OptDirection("asc"),
	)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, errors.New("response is nil")
	}

	vs.total = r.Total

	return videos, nil
}

func (vs *Scraper) GetVideo(id int) (*vimeo.Video, error) {
	video, _, err := vs.api.Videos.Get(id)
	if err != nil {
		return nil, err
	}

	return video, nil
}
