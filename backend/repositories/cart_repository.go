package repositories

import (
	"github.com/cerenkuru/Ecommerce-GoFiber/models"
	"gorm.io/gorm"
)

// sepet işlemleri için repository
type CartRepository interface {
	GetAll(sessionID string) ([]models.CartItem, error)
	GetByProductID(sessionID string, productID uint) (*models.CartItem, error)
	Add(item *models.CartItem) (*models.CartItem, error)
	Update(id uint, quantity int) (*models.CartItem, error)
	Delete(id uint) error
	DeleteAll(sessionID string) error
}

type cartRepository struct {
	db *gorm.DB
}

// yeni bir sepet repository oluşturur
func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetAll(sessionID string) ([]models.CartItem, error) {
	var items []models.CartItem
	if err := r.db.Where("session_id = ?", sessionID).Preload("Product").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *cartRepository) GetByProductID(sessionID string, productID uint) (*models.CartItem, error) {
	var item models.CartItem
	result := r.db.Where("session_id = ? AND product_id = ?", sessionID, productID).Limit(1).Find(&item)
	if err := result.Error; err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &item, nil
}

func (r *cartRepository) Add(item *models.CartItem) (*models.CartItem, error) {
	if err := r.db.Create(item).Error; err != nil {
		return nil, err
	}
	r.db.Preload("Product").First(item, item.ID)
	return item, nil
}

func (r *cartRepository) Update(id uint, quantity int) (*models.CartItem, error) {
	var item models.CartItem
	if err := r.db.Model(&item).Where("id = ?", id).Update("quantity", quantity).Error; err != nil {
		return nil, err
	}
	r.db.Preload("Product").First(&item, id)
	return &item, nil
}

func (r *cartRepository) Delete(id uint) error {
	return r.db.Delete(&models.CartItem{}, id).Error
}

func (r *cartRepository) DeleteAll(sessionID string) error {
	return r.db.Where("session_id = ?", sessionID).Delete(&models.CartItem{}).Error
}
