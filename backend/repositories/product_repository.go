package repositories

import (
	"github.com/cerenkuru/Ecommerce-GoFiber/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAll(category, search string) ([]models.Product, error)
	GetByID(id uint) (*models.Product, error)
	UpdateStock(id uint, stock int) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetAll(category, search string) ([]models.Product, error) {
	var products []models.Product
	query := r.db

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) UpdateStock(id uint, stock int) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Update("stock", stock).Error
}
