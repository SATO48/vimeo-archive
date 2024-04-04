package model

//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen

import (
	"time"

	"github.com/silentsokolov/go-vimeo/v2/vimeo"
)

type Video struct {
	Id uint64 `objectbox:"id(assignable)"`

	URI           string    `json:"uri,omitempty"`
	Name          string    `json:"name,omitempty"`
	Description   string    `json:"description,omitempty"`
	Link          string    `json:"link,omitempty"`
	Duration      int       `json:"duration,omitempty"`
	Width         int       `json:"width,omitempty"`
	Height        int       `json:"height,omitempty"`
	Language      string    `json:"language,omitempty"`
	CreatedTime   time.Time `json:"created_time,omitempty"`
	ModifiedTime  time.Time `json:"modified_time,omitempty"`
	ReleaseTime   time.Time `json:"release_time,omitempty"`
	ContentRating []string  `json:"content_rating,omitempty"`
	License       string    `json:"license,omitempty"`
	Status        string    `json:"status,omitempty"`
	ResourceKey   string    `json:"resource_key,omitempty"`
}

func VideoFromVimeo(v *vimeo.Video) *Video {
	return &Video{
		Id:            uint64(v.GetID()),
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
		ContentRating: v.ContentRating,
		License:       v.License,
		Status:        v.Status,
		ResourceKey:   v.ResourceKey,
	}
}
