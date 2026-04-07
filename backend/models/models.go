package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Image       string  `json:"image"`
	Category    string  `json:"category"`
	Stock       int     `json:"stock"`
}

type CartItem struct {
	gorm.Model
	SessionID string  `json:"session_id" gorm:"index"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int     `json:"quantity"`
}

// CLONE: Product kopyasını yaptık
func (p Product) Clone() Product {
	clone := p
	clone.Model = gorm.Model{}
	return clone
}

// Sayfa boşken db yüklenmesi için prototip ürünler eklenir
func SeedProducts(db *gorm.DB) {
	var count int64
	db.Model(&Product{}).Count(&count)
	if count > 0 {
		return
	}

	// prototipler
	prototypes := []Product{
		{
			Name:        "Sony WH-1000XM5",
			Description: "Endüstri lideri gürültü engelleme özelliğine sahip kablosuz kulaklık. 30 saate kadar pil ömrü, multipoint bağlantı ve Speak-to-Chat teknolojisi.",
			Price:       8499.99,
			Image:       "https://images.unsplash.com/photo-1618366712010-f4ae9c647dcb?w=500",
			Category:    "Elektronik",
			Stock:       15,
		},
		{
			Name:        "Apple MacBook Air M3",
			Description: "M3 çipiyle güçlendirilmiş, 18 saate kadar pil ömrü sunan ince ve hafif laptop. 8GB RAM, 256GB SSD.",
			Price:       42999.00,
			Image:       "https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=500",
			Category:    "Bilgisayar",
			Stock:       8,
		},
		{
			Name:        "Logitech MX Master 3S",
			Description: "8K DPI hassasiyete sahip, sessiz tıklama özellikli, ergonomik kablosuz mouse. Cam yüzeylerle uyumlu.",
			Price:       2299.00,
			Image:       "https://images.unsplash.com/photo-1527864550417-7fd91fc51a46?w=500",
			Category:    "Aksesuar",
			Stock:       30,
		},
		{
			Name:        "Samsung 4K OLED TV 55\"",
			Description: "55 inç 4K OLED panel, 144Hz yenileme hızı, HDR10+ desteği ve Tizen işletim sistemi ile akıllı TV deneyimi.",
			Price:       29999.00,
			Image:       "https://images.unsplash.com/photo-1593359677879-a4bb92f829d1?w=500",
			Category:    "Elektronik",
			Stock:       5,
		},
		{
			Name:        "iPad Pro 12.9\" M4",
			Description: "M4 çipli, Ultra Retina XDR ekranlı profesyonel tablet. Apple Pencil Pro ve Magic Keyboard desteği.",
			Price:       38999.00,
			Image:       "https://images.unsplash.com/photo-1544244015-0df4b3ffc6b0?w=500",
			Category:    "Tablet",
			Stock:       12,
		},
		{
			Name:        "Mechanical Keyboard Keychron Q1",
			Description: "Alüminyum gövdeli, özelleştirilebilir RGB aydınlatmalı, hot-swap mekanik klavye. Gateron G Pro switch.",
			Price:       3799.00,
			Image:       "https://images.unsplash.com/photo-1587829741301-dc798b83add3?w=500",
			Category:    "Aksesuar",
			Stock:       20,
		},
		{
			Name:        "DJI Mini 4 Pro",
			Description: "4K/60fps video çekimi yapabilen, 34 dakika uçuş süreli, obstacle avoidance sensörlü profesyonel drone.",
			Price:       24999.00,
			Image:       "https://images.unsplash.com/photo-1473968512647-3e447244af8f?w=500",
			Category:    "Drone",
			Stock:       7,
		},
		{
			Name:        "Philips Hue Starter Kit",
			Description: "4 adet akıllı ampul ve Hue Bridge içeren başlangıç seti. 16 milyon renk seçeneği, ses ve uygulama kontrolü.",
			Price:       2799.00,
			Image:       "https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=500",
			Category:    "Akıllı Ev",
			Stock:       25,
		},
		{
			Name:        "GoPro Hero 12 Black",
			Description: "5.3K video, 27MP fotoğraf çekimi yapabilen aksiyon kamerası. HyperSmooth 6.0 stabilizasyon, su geçirmez.",
			Price:       13499.00,
			Image:       "https://images.unsplash.com/photo-1516724562728-afc824a36e84?w=500",
			Category:    "Kamera",
			Stock:       18,
		},
		{
			Name:        "Anker PowerBank 26800mAh",
			Description: "26800mAh kapasiteli, 65W hızlı şarj destekli, 3 portlu güç bankası. Laptop şarjı yapabilir.",
			Price:       1899.00,
			Image:       "https://images.unsplash.com/photo-1609091839311-d5365f9ff1c5?w=500",
			Category:    "Aksesuar",
			Stock:       40,
		},
		{
			Name:        "Bose QuietComfort Ultra",
			Description: "Spatial Audio desteği, CustomTune teknolojisi ve 24 saat pil ömrüyle premium gürültü engelleme kulaklık.",
			Price:       11999.00,
			Image:       "https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=500",
			Category:    "Elektronik",
			Stock:       11,
		},
		{
			Name:        "Dell UltraSharp 27\" 4K",
			Description: "27 inç IPS panel, 4K çözünürlük, USB-C 90W şarj, Thunderbolt 4 desteği. Profesyonel renk doğruluğu.",
			Price:       18999.00,
			Image:       "https://images.unsplash.com/photo-1527443224154-c4a3942d3acf?w=500",
			Category:    "Bilgisayar",
			Stock:       9,
		},
	}

	products := make([]Product, 0, len(prototypes))
	for _, prototype := range prototypes {
		products = append(products, prototype.Clone())
	}

	db.Create(&products)
}
