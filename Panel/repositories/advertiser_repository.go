package repositories

import (
	"go-ad-panel/models"
	"gorm.io/gorm"
)

// All functions is Okey

type AdvertiserRepository struct {
	Db *gorm.DB
}

func (t AdvertiserRepository) Save(a *models.Advertiser) error {
	result := t.Db.Create(&a)
	return result.Error
}

func (t AdvertiserRepository) FindByID(id uint) (models.Advertiser, error) {
	var advertiser models.Advertiser
	result := t.Db.First(&advertiser, id)
	return advertiser, result.Error
}

func (t AdvertiserRepository) Update(a *models.Advertiser) error {
	result := t.Db.Save(&a)
	return result.Error
}

func (t AdvertiserRepository) Delete(id uint) error {
	result := t.Db.Delete(&models.Advertiser{}, id)
	return result.Error
}

func (t AdvertiserRepository) FindAll() ([]models.Advertiser, error) {
	var advertisers []models.Advertiser
	result := t.Db.Find(&advertisers)
	return advertisers, result.Error
}
func (t AdvertiserRepository) FindByIDWithAds(id uint) (models.Advertiser, []models.Ad, error) {
	var advertiser models.Advertiser
	var ads []models.Ad
	result := t.Db.First(&advertiser, id)
	if result.Error != nil {
		return advertiser, ads, result.Error
	}
	adsResult := t.Db.Where("advertiser_id = ?", id).Order("title ASC").Find(&ads)
	return advertiser, ads, adsResult.Error
}



var _ AdvertiserRepositoryInterface = (*AdvertiserRepository)(nil)

func (t AdvertiserRepository) FindByIDTx(tx *gorm.DB, id int) (models.Advertiser, error) {
	var advertiser models.Advertiser
	err := tx.First(&advertiser, id).Error
	return advertiser, err
}
func (t AdvertiserRepository) UpdateTx(tx *gorm.DB, advertiser *models.Advertiser) error {
	return tx.Save(advertiser).Error
}
func (t AdvertiserRepository) DecreaseCredit(tx *gorm.DB, advertiser *models.Advertiser, bid int) error {
	return tx.Model(advertiser).Update("Credit", gorm.Expr("Credit - ?", bid)).Error
}
