package app

import (
	"github.com/defval/di"
	"github.com/sato48/vimeo-archive/lib/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Gorm(c *AppContainer) error {
	return c.Apply(
		di.Provide(func() (*gorm.DB, error) {
			db, err := gorm.Open(sqlite.Open("vimeo-archive.db"), &gorm.Config{})
			if err == nil {
				db.AutoMigrate(&model.Video{})
				db.AutoMigrate(&model.File{})
			}
			return db, err
		}),
	)
}
