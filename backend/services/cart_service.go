package services

import (
	"errors"
	"math"

	"github.com/cerenkuru/Ecommerce-GoFiber/models"
	"github.com/cerenkuru/Ecommerce-GoFiber/repositories"
)

const kdvRate = 0.20

type CartSummaryLine struct {
	CartItemID uint    `json:"cart_item_id"`
	ProductID  uint    `json:"product_id"`
	Product    string  `json:"product"`
	KDV        float64 `json:"kdv"`
}

type CartSummary struct {
	Subtotal float64           `json:"subtotal"`
	KDVRate  float64           `json:"kdv_rate"`
	KDV      float64           `json:"kdv"`
	Total    float64           `json:"total"`
	Lines    []CartSummaryLine `json:"lines"`
}

// CartService defines cart business logic
type CartService interface {
	GetCart(sessionID string) ([]models.CartItem, error)
	GetCartSummary(sessionID string) (*CartSummary, error)
	AddToCart(sessionID string, productID uint, quantity int) (*models.CartItem, error)
	UpdateCartItem(sessionID string, id uint, quantity int) (*models.CartItem, error)
	RemoveFromCart(id uint) error
	ClearCart(sessionID string) error
	Checkout(sessionID string) (string, error)
}

type cartService struct {
	repo        repositories.CartRepository
	productRepo repositories.ProductRepository
}

// NewCartService creates a new cart service
func NewCartService(repo repositories.CartRepository, productRepo repositories.ProductRepository) CartService {
	return &cartService{repo: repo, productRepo: productRepo}
}

func (s *cartService) GetCart(sessionID string) ([]models.CartItem, error) {
	return s.repo.GetAll(sessionID)
}

func (s *cartService) GetCartSummary(sessionID string) (*CartSummary, error) {
	items, err := s.repo.GetAll(sessionID)
	if err != nil {
		return nil, err
	}

	var subtotal float64
	var kdvTotal float64
	lines := make([]CartSummaryLine, 0, len(items))

	for _, item := range items {
		lineSubtotal := item.Product.Price * float64(item.Quantity)
		subtotal += lineSubtotal

		lineKDV := round2(lineSubtotal * kdvRate)
		kdvTotal += lineKDV

		lines = append(lines, CartSummaryLine{
			CartItemID: item.ID,
			ProductID:  item.ProductID,
			Product:    item.Product.Name,
			KDV:        lineKDV,
		})
	}

	subtotal = round2(subtotal)
	kdvTotal = round2(kdvTotal)
	total := round2(subtotal + kdvTotal)

	return &CartSummary{
		Subtotal: subtotal,
		KDVRate:  kdvRate,
		KDV:      kdvTotal,
		Total:    total,
		Lines:    lines,
	}, nil
}

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}

func (s *cartService) AddToCart(sessionID string, productID uint, quantity int) (*models.CartItem, error) {
	if quantity <= 0 {
		quantity = 1
	}

	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, errors.New("ürün bulunamadı")
	}

	if product.Stock <= 0 {
		return nil, errors.New("urun stokta yok")
	}

	if product.Stock < quantity {
		return nil, errors.New("stok yetersiz")
	}

	existing, err := s.repo.GetByProductID(sessionID, productID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		newQuantity := existing.Quantity + quantity
		if product.Stock < newQuantity {
			return nil, errors.New("stok yetersiz")
		}
		return s.repo.Update(existing.ID, newQuantity)
	}

	item := &models.CartItem{
		SessionID: sessionID,
		ProductID: productID,
		Quantity:  quantity,
	}
	return s.repo.Add(item)
}

func (s *cartService) UpdateCartItem(sessionID string, id uint, quantity int) (*models.CartItem, error) {
	if quantity <= 0 {
		return nil, errors.New("geçersiz miktar")
	}

	items, err := s.repo.GetAll(sessionID)
	if err != nil {
		return nil, err
	}

	var target *models.CartItem
	for i := range items {
		if items[i].ID == id {
			target = &items[i]
			break
		}
	}

	if target == nil {
		return nil, errors.New("sepet öğesi bulunamadı")
	}

	product, err := s.productRepo.GetByID(target.ProductID)
	if err != nil {
		return nil, errors.New("ürün bulunamadı")
	}

	if quantity > product.Stock {
		return nil, errors.New("stok yetersiz")
	}

	return s.repo.Update(id, quantity)
}

func (s *cartService) RemoveFromCart(id uint) error {
	return s.repo.Delete(id)
}

func (s *cartService) ClearCart(sessionID string) error {
	return s.repo.DeleteAll(sessionID)
}

func (s *cartService) Checkout(sessionID string) (string, error) {
	items, err := s.repo.GetAll(sessionID)
	if err != nil {
		return "", err
	}

	if len(items) == 0 {
		return "", errors.New("sepet bos")
	}

	products := make(map[uint]*models.Product)
	for _, item := range items {
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			return "", errors.New("urun bulunamadi")
		}

		if product.Stock < item.Quantity {
			return "", errors.New(product.Name + " icin stok yetersiz")
		}

		products[item.ProductID] = product
	}

	for _, item := range items {
		product := products[item.ProductID]
		newStock := product.Stock - item.Quantity
		if err := s.productRepo.UpdateStock(item.ProductID, newStock); err != nil {
			return "", err
		}
	}

	if err := s.repo.DeleteAll(sessionID); err != nil {
		return "", err
	}

	return "Siparis olusturuldu", nil
}
