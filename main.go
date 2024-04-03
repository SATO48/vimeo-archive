package main

import (
	"log/slog"
	"os"
	"vimeo-archive/lib/vimeo"

	"github.com/spf13/viper"
)

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
}

func main() {
	slog.Info("identified archived videos", "count", len(archived))

	vs := vimeo.NewScraper(
		vimeo.WithPageSize(25),
		vimeo.WithMaxPages(2),
	)

	for vs.HasNextPage() {
		videos, err := vs.ListVideos()
		if err != nil {
			slog.Error("unable to list videos", "error", err)
			os.Exit(1)
		}

		slog.Debug("listed vimeo videos", "count", len(videos))

		for _, v := range videos {
			if isArchived(v.GetID()) {
				slog.Debug("found archived video", "id", v.GetID())
				continue
			}

			slog.Warn("found unarchived video", "id", v.GetID(), "name", v.Name)

			f := vimeo.FindBestFile(v.Files)
			slog.Info("found best video file", "link", f.Link[:12], "quality", f.Quality, "h", f.Height, "w", f.Width)
		}
	}
}

func init() {
	viper.AutomaticEnv()
}
