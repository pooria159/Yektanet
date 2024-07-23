package models

import "gorm.io/gorm"

type Publisher struct {
	gorm.Model
	Name    string `gorm:"type:varchar(255)"`
	Website string `gorm:"type:varchar(255)"`
	Credit  int    `gorm:"type:int"`
}
