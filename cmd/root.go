package cmd

import (
	"log/slog"
	"vimeo-archive/app"
	"vimeo-archive/lib/model"
	libvimeo "vimeo-archive/lib/vimeo"

	"github.com/defval/di"
	"github.com/silentsokolov/go-vimeo/v2/vimeo"

	"github.com/spf13/cobra"
)

func Bootstrap(c *app.AppContainer) error {
	rootCmd := &cobra.Command{
		Use:   "vimeo-archiver",
		Short: "Vimeo Archiver is a tool to archive Vimeo videos",
		RunE: c.RunE(func(vc *vimeo.Client, vb *model.VideoBox) error {
			vs := libvimeo.NewScraper(
				libvimeo.WithAPI(vc),
				libvimeo.WithPageSize(25),
				libvimeo.WithPagePointer(88),
				libvimeo.WithMaxPages(2),
			)

			for vs.HasNextPage() {
				videos, err := vs.ListVideos()
				if err != nil {
					slog.Error("unable to list videos", "error", err)
					return err
				}

				slog.Debug("listed vimeo videos", "count", len(videos))

				for _, v := range videos {
					video := model.VideoFromVimeo(v)
					_, err := vb.Put(video)
					if err != nil {
						slog.Error("unable to put video", "error", err)
						return err
					}
				}
			}

			return nil
		}),
	}

	return c.Apply(
		di.ProvideValue(rootCmd),
	)
}
