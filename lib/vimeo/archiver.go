package vimeo

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sato48/vimeo-archive/lib/model"
	"github.com/silentsokolov/go-vimeo/v2/vimeo"
	"gorm.io/gorm"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

type Archiver struct {
	s3  *s3.Client
	up  *manager.Uploader
	vc  *vimeo.Client
	db  *gorm.DB
	max uint64
}

type ArchiverOptionFunc func(*Archiver)

func WithS3Client(s3 *s3.Client) ArchiverOptionFunc {
	return func(a *Archiver) {
		a.s3 = s3
		a.up = manager.NewUploader(s3)
	}
}

func WithVimeoClient(vc *vimeo.Client) ArchiverOptionFunc {
	return func(a *Archiver) {
		a.vc = vc
	}
}

func WithDB(db *gorm.DB) ArchiverOptionFunc {
	return func(a *Archiver) {
		a.db = db
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
	videos := []*model.Video{}
	if err := a.db.Model(&model.Video{}).
		Where("downloaded_time IS NULL").
		Limit(int(a.max)).
		Find(&videos).
		Error; err != nil {
		return err
	}

	slog.Info("archiving videos", "count", len(videos))

	for _, video := range videos {
		if err := a.archiveVideo(video); err != nil {
			return err
		}

		slog.Info("archived video", "id", video.ID)
	}

	return nil
}

func (a *Archiver) archiveVideo(v *model.Video) error {
	r, _, err := a.vc.Videos.Get(int(v.ID))
	if err != nil {
		return err
	}
	f := FindBestFile(r.Files)

	// Download the file
	dl, err := http.Get(f.Link)
	if err != nil {
		slog.Error("unable to download file", "error", err)
		return err
	}
	defer dl.Body.Close()

	// Upload the file
	dst := fmt.Sprintf("%d%s", v.ID, path.Ext(dl.Request.URL.Path))
	_, err = a.up.Upload(context.Background(), &s3.PutObjectInput{
		Bucket:             aws.String("sato48-vimeo"),
		Key:                aws.String(dst),
		Body:               dl.Body,
		ContentType:        aws.String(dl.Header.Get("Content-Type")),
		ContentDisposition: aws.String("inline"),
	})
	if err != nil {
		slog.Error("unable to upload file", "error", err)
		return err
	}

	// Insert File and mark as Downloaded
	fm := model.FileFromVimeo(f)
	fm.VideoID = v.ID
	if err := a.db.Create(&fm).Error; err != nil {
		slog.Error("unable to put file", "error", err)
		return err
	}

	now := time.Now()
	v.DownloadedTime = &now
	if err := a.db.Save(v).Error; err != nil {
		slog.Error("unable to update video", "error", err)
		return err
	}

	return nil
}
