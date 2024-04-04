package app

import (
	"vimeo-archive/lib/model"

	"github.com/defval/di"
	"github.com/objectbox/objectbox-go/objectbox"
)

func Objectbox(c *AppContainer) error {
	return c.Apply(
		di.Provide(func() (*objectbox.ObjectBox, error) {
			return objectbox.NewBuilder().Model(model.ObjectBoxModel()).Build()
		}),
		di.Provide(func(ob *objectbox.ObjectBox) *model.VideoBox {
			return model.BoxForVideo(ob)
		}),
		di.Provide(func(ob *objectbox.ObjectBox) *model.FileBox {
			return model.BoxForFile(ob)
		}),
	)
}
