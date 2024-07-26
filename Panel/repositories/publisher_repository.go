package repositories

import (
	"go-ad-panel/models"
	"gorm.io/gorm"
)

// All Functions is okey

type PublisherRepository struct {
	Db *gorm.DB
}

func (t PublisherRepository) Save(p *models.Publisher) error {
	result := t.Db.Create(&p)
	return result.Error
}

func (t PublisherRepository) FindByID(id uint) (models.Publisher, error) {
	var publisher models.Publisher
	result := t.Db.First(&publisher, id)
	return publisher, result.Error
}

var _ PublisherRepositoryInterface = (*PublisherRepository)(nil)

func (t PublisherRepository) Update(p *models.Publisher) error {
	result := t.Db.Save(&p)
	return result.Error
}

func (t PublisherRepository) Delete(id uint) error {
	result := t.Db.Delete(&models.Publisher{}, id)
	return result.Error
}

func (t PublisherRepository) FindAll() ([]models.Publisher, error) {
	var publishers []models.Publisher
	result := t.Db.Find(&publishers)
	return publishers, result.Error
}

var _ PublisherRepositoryInterface = (*PublisherRepository)(nil)

func (t PublisherRepository) FindByIDTx(tx *gorm.DB, id int) (models.Publisher, error) {
	var publisher models.Publisher
	err := tx.First(&publisher, id).Error
	return publisher, err
}
func (t PublisherRepository) UpdateTx(tx *gorm.DB, publisher *models.Publisher) error {
	return tx.Save(publisher).Error
}
func (t PublisherRepository) IncreaseCredit(tx *gorm.DB, publisher *models.Publisher, bid int) error {
	return tx.Model(publisher).Update("Credit", gorm.Expr("Credit + ?", bid)).Error
}
