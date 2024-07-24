package repositories

import "go-ad-panel/models"


type PublisherRepositoryInterface interface {
    FindByID(id uint) (models.Publisher, error)
    Save(p *models.Publisher) error
    Update(p *models.Publisher) error
    Delete(id uint) error
    FindAll() ([]models.Publisher, error)
}
