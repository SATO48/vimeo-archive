package model

//go:generate go run github.com/objectbox/objectbox-go/cmd/objectbox-gogen

import (
	"time"

	"github.com/silentsokolov/go-vimeo/v2/vimeo"
)

type File struct {
	Id    uint64
	Video *Video `objectbox:"link"`

	FileID      string    `json:"video_file_id,omitempty" objectbox:"unique"`
	Quality     string    `json:"quality,omitempty"`
	Type        string    `json:"type,omitempty"`
	Width       int       `json:"width,omitempty"`
	Height      int       `json:"height,omitempty"`
	Link        string    `json:"link,omitempty"`
	CreatedTime time.Time `json:"created_time,omitempty"`
	FPS         float32   `json:"fps,omitempty"`
	Size        int       `json:"size,omitempty"`
	MD5         string    `json:"md5,omitempty"`
}

func FileFromVimeo(f *vimeo.File) *File {
	return &File{
		FileID:      f.FileID,
		Quality:     f.Quality,
		Type:        f.Type,
		Width:       f.Width,
		Height:      f.Height,
		Link:        f.Link,
		CreatedTime: f.CreatedTime,
		FPS:         f.FPS,
		Size:        f.Size,
		MD5:         f.MD5,
	}
}
