package signature

import (
	"gorm.io/gorm"
)

func AutoOrderedMigrate(db *gorm.DB, values ...interface{}) error {
	for _, value := range values {
		if !db.Migrator().HasTable(value) {
			if err := db.Migrator().CreateTable(value); err != nil {
				println()
				return err
			}
		}
	}

	return nil
}
