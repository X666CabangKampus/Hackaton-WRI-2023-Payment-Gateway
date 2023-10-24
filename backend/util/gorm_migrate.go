package signature

import (
	"gorm.io/gorm"
)

func AutoOrderedMigrate(db *gorm.DB, values ...interface{}) error {
	for _, value := range values {
		execTx := db
		if !db.Migrator().HasTable(value) {
			if err := execTx.Migrator().CreateTable(value); err != nil {
				return err
			}
		}
	}

	return nil
}
