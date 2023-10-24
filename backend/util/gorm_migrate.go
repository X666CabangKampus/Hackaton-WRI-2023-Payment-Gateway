package signature

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
)

func AutoOrderedMigrate(db *gorm.DB, values ...interface{}) error {
	for _, value := range values {
		if !db.Migrator().HasTable(value) {
			if err := db.Migrator().CreateTable(value); err != nil {
				fmt.Printf("ERR: %v\n", reflect.TypeOf(value))
				return err
			}
		}
	}

	return nil
}
