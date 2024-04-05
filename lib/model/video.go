package model

import (
	"strings"
	"time"

	"github.com/silentsokolov/go-vimeo/v2/vimeo"
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Files []File

	URI            string
	Name           string
	Description    string
	Link           string
	Duration       int
	Width          int
	Height         int
	Language       string
	CreatedTime    time.Time
	ModifiedTime   time.Time
	ReleaseTime    time.Time
	ContentRating  string
	License        string
	Status         string
	ResourceKey    string
	DownloadedTime *time.Time
}

func VideoFromVimeo(v *vimeo.Video) *Video {
	return &Video{
		Model: gorm.Model{
			ID: uint(v.GetID()),
		},
		URI:           v.URI,
		Name:          v.Name,
		Description:   v.Description,
		Link:          v.Link,
		Duration:      v.Duration,
		Width:         v.Width,
		Height:        v.Height,
		Language:      v.Language,
		CreatedTime:   v.CreatedTime,
		ModifiedTime:  v.ModifiedTime,
		ReleaseTime:   v.ReleaseTime,
		ContentRating: strings.Join(v.ContentRating, ","),
		License:       v.License,
		Status:        v.Status,
		ResourceKey:   v.ResourceKey,
	}
}
