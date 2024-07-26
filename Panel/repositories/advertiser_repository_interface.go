package repositories

import "go-ad-panel/models"




type AdvertiserRepositoryInterface interface {
    FindByID(id uint) (models.Advertiser, error)
    Save(a *models.Advertiser) error
    Update(a *models.Advertiser) error
    Delete(id uint) error
    FindAll() ([]models.Advertiser, error)
    FindByIDWithAds(id uint) (models.Advertiser, []models.Ad, error)
}
