package models

type Publisher struct {
	ID     int    `gorm:"type:int;primary_key"`
	Name    string `gorm:"type:varchar(255)"`
	Website string `gorm:"type:varchar(255)"`
	Credit  int    `gorm:"type:int"`
}
