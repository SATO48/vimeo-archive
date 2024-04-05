package model

import (
	"time"

	"github.com/silentsokolov/go-vimeo/v2/vimeo"
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	VideoID uint
	Video   Video

	FileID      string    `db:"video_file_id,omitempty"`
	Quality     string    `db:"quality,omitempty"`
	Type        string    `db:"type,omitempty"`
	Width       int       `db:"width,omitempty"`
	Height      int       `db:"height,omitempty"`
	Link        string    `db:"link,omitempty"`
	CreatedTime time.Time `db:"created_time,omitempty"`
	FPS         float32   `db:"fps,omitempty"`
	Size        int       `db:"size,omitempty"`
	MD5         string    `db:"md5,omitempty"`
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
