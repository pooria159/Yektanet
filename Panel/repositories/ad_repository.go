package repositories

import (
	"go-ad-panel/models"
	"gorm.io/gorm"
)

type AdRepository struct {
	Db *gorm.DB
}

func (t AdRepository) WithTransaction(fn func(tx *gorm.DB) error) error {
	return t.Db.Transaction(fn)
}

func (t AdRepository) FindByIDTx(tx *gorm.DB, id int) (models.Ad, error) {
	var ad models.Ad
	err := tx.First(&ad, id).Error
	return ad, err
}
func (t AdRepository) UpdateTx(tx *gorm.DB, ad *models.Ad) error {
	return tx.Save(ad).Error
}
func (t AdRepository) IncrementImpressionsTx(tx *gorm.DB, ad *models.Ad) error {
	return tx.Model(ad).Update("Impressions", gorm.Expr("Impressions + ?", 1)).Error
}

//	func (t AdRepository) IncrementClicksTx(tx *gorm.DB, ad *models.Ad) error {
//		return tx.Model(ad).Update("Clicks", gorm.Expr("Clicks + ?", 1)).Error
//	}
func (t AdRepository) IncrementClicksTx(tx *gorm.DB, ad *models.Ad) error {
	return tx.Model(ad).Updates(map[string]interface{}{
		"Clicks":        gorm.Expr("Clicks + ?", 1),
		"EngagedCredit": gorm.Expr("engaged_credit + ?", ad.BidValue),
	}).Error
}

func (t AdRepository) Save(ad models.Ad) error {
	result := t.Db.Create(&ad)
	return result.Error
}

func (t AdRepository) FindByID(id int) (models.Ad, error) {
	var ad models.Ad
	result := t.Db.First(&ad, id)
	return ad, result.Error
}

func (t AdRepository) Update(ad *models.Ad) error {
	result := t.Db.Save(&ad)
	return result.Error
}

func (t AdRepository) Delete(id int) error {
	result := t.Db.Delete(&models.Ad{}, id)
	return result.Error
}

func (t AdRepository) FindAll() ([]models.Ad, error) {
	var ads []models.Ad
	result := t.Db.Find(&ads)
	return ads, result.Error
}

func (t AdRepository) FindAllAdsByAdvertiser(advertiserID int) ([]models.Ad, error) {
	var ads []models.Ad
	result := t.Db.Where("advertiser_id = ?", advertiserID).Find(&ads)
	return ads, result.Error
}

func (t AdRepository) FindAllActiveAds() ([]models.Ad, error) {
	var ads []models.Ad
	result := t.Db.Where("is_active = ?", true).Find(&ads)
	return ads, result.Error
}
