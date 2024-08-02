package models

import "gorm.io/gorm"

type Advertiser struct {
	gorm.Model
	Name   string `gorm:"type:varchar(255)"`
	Credit int    `gorm:"type:int"`
	Ads    []Ad
}
