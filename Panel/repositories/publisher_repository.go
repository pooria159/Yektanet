package repositories

import (
	"go-ad-panel/models"
)

type TagsRepository interface {
	Save(tags models.Publisher)
	Update(tags models.Publisher)
	Delete(PublisherId int)
	FindById(PublisherId int) (publishers models.Publisher, err error)
	FindAll() []models.Publisher
}
