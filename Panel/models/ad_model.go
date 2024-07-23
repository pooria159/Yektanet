package models

import "gorm.io/gorm"

type Ad struct {
	gorm.Model
	Title        string     `gorm:"type:varchar(255);not null"`
	ImagePath    string     `gorm:"type:varchar(255);not null"`
	BidValue     int        `gorm:"type:int;not null"`
	IsActive     bool       `gorm:"type:boolean;not null"`
	AdvertiserID int        `gorm:"type:int;not null"`
	Advertiser   Advertiser `gorm:"foreignKey:AdvertiserID"`
}
