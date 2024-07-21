package models

type Publisher struct {
	Id      int    `gorm:"type:int;primary_key"`
	Name    string `gorm:"type:varchar(255)"`
	Website string `gorm:"type:varchar(255)"`
	credit  int    `gorm:"type:int"`
}
