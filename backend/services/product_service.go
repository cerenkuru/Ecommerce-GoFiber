package services

import (
	"errors"

	"github.com/cerenkuru/Ecommerce-GoFiber/models"
	"github.com/cerenkuru/Ecommerce-GoFiber/repositories"
)

// ProductService defines product business logic
type ProductService interface {
	GetAllProducts(category, search string) ([]models.Product, error)
	GetProductByID(id uint) (*models.Product, error)
}

type productService struct {
	repo repositories.ProductRepository
}

// NewProductService creates a new product service
func NewProductService(repo repositories.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetAllProducts(category, search string) ([]models.Product, error) {
	return s.repo.GetAll(category, search)
}

func (s *productService) GetProductByID(id uint) (*models.Product, error) {
	if id == 0 {
		return nil, errors.New("geçersiz ürün ID")
	}

	product, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("ürün bulunamadı")
	}
	return product, nil
}
