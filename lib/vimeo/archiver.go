package vimeo

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"path"
	"regexp"
	"time"
	"vimeo-archive/lib/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/silentsokolov/go-vimeo/v2/vimeo"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

type Archiver struct {
	s3  *s3.Client
	up  *manager.Uploader
	vc  *vimeo.Client
	vb  *model.VideoBox
	fb  *model.FileBox
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

func WithVideoBox(vb *model.VideoBox) ArchiverOptionFunc {
	return func(a *Archiver) {
		a.vb = vb
	}
}

func WithFileBox(fb *model.FileBox) ArchiverOptionFunc {
	return func(a *Archiver) {
		a.fb = fb
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

	slog.Info("archiving videos", "count", len(videos))

	for _, video := range videos {
		if err := a.archiveVideo(video); err != nil {
			return err
		}

		slog.Info("archived video", "id", video.Id)
	}

	return nil
}

func (a *Archiver) archiveVideo(v *model.Video) error {
	r, _, err := a.vc.Videos.Get(int(v.Id))
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
	dst := fmt.Sprintf("%d%s", v.Id, path.Ext(dl.Request.URL.Path))
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

	// TODO: Insert File and mark as Downloaded
	fm := model.FileFromVimeo(f)
	fm.Video = v
	if _, err = a.fb.Put(fm); err != nil {
		slog.Error("unable to put file", "error", err)
		return err
	}

	v.DownloadedTime = time.Now()
	if _, err = a.vb.Put(v); err != nil {
		slog.Error("unable to put video", "error", err)
		return err
	}

	return nil
}
