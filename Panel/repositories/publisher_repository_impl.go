package repositories

import (
	"go-ad-panel/models"
	"gorm.io/gorm"
)

type PublisherRepositoryImpl struct {
	Db *gorm.DB
}

func (t PublisherRepositoryImpl) Save(p models.Publisher) error {
	result := t.Db.Create(&p)
	return result.Error
}

func (t PublisherRepositoryImpl) FindByID(id uint) (models.Publisher, error) {
	var publisher models.Publisher
	result := t.Db.First(&publisher, id)
	return publisher, result.Error
}

func (t PublisherRepositoryImpl) Update(p models.Publisher) error {
	result := t.Db.Save(&p)
	return result.Error
}

func (t PublisherRepositoryImpl) Delete(id uint) error {
	result := t.Db.Delete(&models.Publisher{}, id)
	return result.Error
}

func (t PublisherRepositoryImpl) FindAll() ([]models.Publisher, error) {
	var publishers []models.Publisher
	result := t.Db.Find(&publishers)
	return publishers, result.Error
}
