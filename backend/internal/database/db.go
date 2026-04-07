package database

import (
	"fmt"
	"os"
	"sync"

	"github.com/cerenkuru/Ecommerce-GoFiber/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	once       sync.Once // bu kod 1 kere çalışsın
	initErr    error
)

// Singleton: bir kere db bağlantısı kurulur
func GetDB() (*gorm.DB, error) {
	once.Do(func() {
		connectionStr := os.Getenv("DB_USER") + ":" +
			os.Getenv("DB_PASS") + "@tcp(" +
			os.Getenv("DB_HOST") + ":" +
			os.Getenv("DB_PORT") + ")/" +
			os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"

		db, err := gorm.Open(mysql.Open(connectionStr), &gorm.Config{})
		if err != nil {
			initErr = fmt.Errorf("MySQL baglanti hatasi: %w", err)
			return
		}

		if err := db.AutoMigrate(&models.Product{}, &models.CartItem{}); err != nil {
			initErr = fmt.Errorf("migrate hatasi: %w", err)
			return
		}

		db.Exec("ALTER TABLE products DROP COLUMN rating;")

		models.SeedProducts(db)
		dbInstance = db

	})

	return dbInstance, initErr
}
