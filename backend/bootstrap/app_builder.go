package bootstrap

import (
	"errors"

	"github.com/cerenkuru/Ecommerce-GoFiber/handlers"
	"github.com/cerenkuru/Ecommerce-GoFiber/internal/database"
	"github.com/cerenkuru/Ecommerce-GoFiber/repositories"
	"github.com/cerenkuru/Ecommerce-GoFiber/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

type AppBuilder struct {
	app     *fiber.App
	db      *gorm.DB
	handler *handlers.Handler
	err     error
}

func NewAppBuilder() *AppBuilder {
	return &AppBuilder{app: fiber.New()}
}

func (b *AppBuilder) WithDatabase() *AppBuilder {
	if b.err != nil {
		return b
	}

	db, err := database.GetDB()
	if err != nil {
		b.err = err
		return b
	}

	b.db = db
	return b
}

func (b *AppBuilder) WithDependencies() *AppBuilder {
	if b.err != nil {
		return b
	}

	if b.db == nil {
		b.err = errors.New("database baslatilmadan dependency kurulamaz")
		return b
	}

	productRepo := repositories.NewProductRepository(b.db)
	cartRepo := repositories.NewCartRepository(b.db)
	productService := services.NewProductService(productRepo)
	cartService := services.NewCartService(cartRepo, productRepo)
	b.handler = handlers.NewHandler(productService, cartService)

	return b
}

func (b *AppBuilder) WithMiddleware() *AppBuilder {
	if b.err != nil {
		return b
	}

	b.app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))
	return b
}

func (b *AppBuilder) WithRoutes() *AppBuilder {
	if b.err != nil {
		return b
	}

	if b.handler == nil {
		b.err = errors.New("handler olusmadi")
		return b
	}

	api := b.app.Group("/api")
	api.Get("/products", b.handler.GetProducts)
	api.Get("/products/:id", b.handler.GetProduct)
	api.Get("/cart", b.handler.GetCart)
	api.Get("/cart/summary", b.handler.GetCartSummary)
	api.Get("/cart/details", b.handler.GetCartDetails)
	api.Post("/cart", b.handler.AddToCart)
	api.Post("/cart/checkout", b.handler.Checkout)
	api.Put("/cart/:id", b.handler.UpdateCartItem)
	api.Delete("/cart/:id", b.handler.DeleteCartItem)
	api.Delete("/cart", b.handler.ClearCart)
	return b
}

func (b *AppBuilder) Build() (*fiber.App, error) {
	if b.err != nil {
		return nil, b.err
	}

	if b.db == nil {
		return nil, errors.New("database baslatilamadi")
	}
	if b.handler == nil {
		return nil, errors.New("bagimliliklar baslatilamadi")
	}
	return b.app, nil
}
