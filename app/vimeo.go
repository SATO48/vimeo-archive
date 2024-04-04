package app

import (
	"context"

	"github.com/silentsokolov/go-vimeo/v2/vimeo"

	"github.com/defval/di"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func Vimeo(c *AppContainer) error {
	return c.Apply(
		di.Provide(func() *vimeo.Client {
			return vimeo.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: viper.GetString("VIMEO_ACCESS_KEY")},
			)), nil)
		}),
	)
}
