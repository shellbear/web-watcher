package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

func New() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "db.sqlite")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to models: %s", err)
	}

	if err = db.AutoMigrate(&Task{}).Error; err != nil {
		return nil, err
	}

	return db, nil
}
