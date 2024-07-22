package models

type Advertiser struct {
	Id     int    `gorm:"type:int;primary_key"`
	Name   string `gorm:"type:varchar(255)"`
	Credit int    `gorm:"type:int"`
}
