package handlers

import (
	"crypto/rand"
	"fmt"
	"strconv"

	"github.com/cerenkuru/Ecommerce-GoFiber/models"
	"github.com/cerenkuru/Ecommerce-GoFiber/services"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	productService services.ProductService
	cartService    services.CartService
}

// Dışarıdan servisleri alarak handler'ı oluşturuyoruz. Böylece handler sadece HTTP ile ilgilenir, iş mantığı servislerde kalır.

func NewHandler(productService services.ProductService, cartService services.CartService) *Handler {
	return &Handler{
		productService: productService,
		cartService:    cartService,
	}
}

// getSessionID fonksiyonu kullanıcıya özel bir session ID oluşturur veya mevcutsa onu döner.
// Bu ID, kullanıcının sepetini tanımlamak için kullanılır.
func (h *Handler) getSessionID(c *fiber.Ctx) string {
	sessionID := c.Cookies("session_id")
	if sessionID == "" {
		sessionID = generateSessionID()
		c.Cookie(&fiber.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			MaxAge:   86400 * 7, // 7 days
			HTTPOnly: true,
		})
	}
	return sessionID
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (h *Handler) GetProducts(c *fiber.Ctx) error {
	category := c.Query("category")
	search := c.Query("search")

	products, err := h.productService.GetAllProducts(category, search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ürünler alınamadı"})
	}

	if len(products) == 0 {
		return c.JSON([]interface{}{})
	}
	return c.JSON(products)
}

func (h *Handler) GetProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz ürün ID"})
	}

	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(product)
}

func (h *Handler) GetCart(c *fiber.Ctx) error {
	sessionID := h.getSessionID(c)
	items, err := h.cartService.GetCart(sessionID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Sepet alınamadı"})
	}

	if len(items) == 0 {
		return c.JSON([]interface{}{})
	}
	return c.JSON(items)
}

func (h *Handler) GetCartSummary(c *fiber.Ctx) error {
	sessionID := h.getSessionID(c)
	summary, err := h.cartService.GetCartSummary(sessionID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Sepet ozeti alınamadı"})
	}

	return c.JSON(summary)
}

func (h *Handler) GetCartDetails(c *fiber.Ctx) error {
	sessionID := h.getSessionID(c)
	items, err := h.cartService.GetCart(sessionID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Sepet alınamadı"})
	}

	summary, err := h.cartService.GetCartSummary(sessionID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Sepet ozeti alınamadı"})
	}

	if len(items) == 0 {
		items = []models.CartItem{}
	}

	return c.JSON(fiber.Map{
		"items":   items,
		"summary": summary,
	})
}

func (h *Handler) AddToCart(c *fiber.Ctx) error {
	sessionID := h.getSessionID(c)
	type Body struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	}

	var body Body
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz istek"})
	}

	item, err := h.cartService.AddToCart(sessionID, body.ProductID, body.Quantity)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(item)
}

func (h *Handler) UpdateCartItem(c *fiber.Ctx) error {
	sessionID := h.getSessionID(c)
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz ID"})
	}

	type Body struct {
		Quantity int `json:"quantity"`
	}

	var body Body
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz istek"})
	}

	if body.Quantity <= 0 {
		if err := h.cartService.RemoveFromCart(uint(id)); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Ürün kaldırılamadı"})
		}
		return c.JSON(fiber.Map{"message": "Ürün sepetten kaldırıldı"})
	}

	item, err := h.cartService.UpdateCartItem(sessionID, uint(id), body.Quantity)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Güncellenemedi"})
	}

	return c.JSON(item)
}

func (h *Handler) DeleteCartItem(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Geçersiz ID"})
	}

	if err := h.cartService.RemoveFromCart(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ürün kaldırılamadı"})
	}

	return c.JSON(fiber.Map{"message": "Ürün sepetten kaldırıldı"})
}

func (h *Handler) ClearCart(c *fiber.Ctx) error {
	sessionID := h.getSessionID(c)
	if err := h.cartService.ClearCart(sessionID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Sepet temizlenemedi"})
	}

	return c.JSON(fiber.Map{"message": "Sepet temizlendi"})
}

func (h *Handler) Checkout(c *fiber.Ctx) error {
	sessionID := h.getSessionID(c)
	message, err := h.cartService.Checkout(sessionID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": message})
}
