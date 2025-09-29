package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// Product represents the product model
type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Price       float64   `json:"price" gorm:"not null"`
	Stock       int       `json:"stock" gorm:"default:0"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=100"`
	Description string  `json:"description" binding:"max=500"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Stock       int     `json:"stock" binding:"min=0"`
	Category    string  `json:"category" binding:"required,min=1,max=50"`
}

type ProductService struct {
	db    *gorm.DB
	Redis *redis.Client
}

// Database connection
func NewProductService() *ProductService {
	// Database connection
	db, err := gorm.Open(sqlserver.Open("sqlserver://lek:P@ssw0rd@10.100.33.68:1433?database=product_db_demo"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate
	db.AutoMigrate(&Product{})

	// Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &ProductService{
		db:    db,
		Redis: rdb,
	}
}

func (ps *ProductService) CreateProduct(c *gin.Context) {
	var productReq ProductRequest
	if err := c.ShouldBindJSON(&productReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := Product{
		Name:        productReq.Name,
		Description: productReq.Description,
		Price:       productReq.Price,
		Stock:       productReq.Stock,
		Category:    productReq.Category,
	}

	if err := ps.db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Cache the product
	productJSON, _ := json.Marshal(product)
	ps.Redis.Set(context.Background(), "product:"+strconv.Itoa(int(product.ID)), productJSON, time.Hour)

	c.JSON(http.StatusCreated, product)
}

func (ps *ProductService) GetProduct(c *gin.Context) {
	id := c.Param("id")

	// Try to get from cache first
	cachedProduct, err := ps.Redis.Get(context.Background(), "product:"+id).Result()
	if err == nil {
		var product Product
		json.Unmarshal([]byte(cachedProduct), &product)
		c.JSON(http.StatusOK, product)
		return
	}

	// If not in cache, get from database
	var product Product
	if err := ps.db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Cache the product
	productJSON, _ := json.Marshal(product)
	ps.Redis.Set(context.Background(), "product:"+id, productJSON, time.Hour)

	c.JSON(http.StatusOK, product)
}

func (ps *ProductService) GetAllProducts(c *gin.Context) {
	var products []Product

	category := c.Query("category")
	query := ps.db

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (ps *ProductService) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var product Product
	if err := ps.db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ps.db.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update cache
	productJSON, _ := json.Marshal(product)
	ps.Redis.Set(context.Background(), "product:"+id, productJSON, time.Hour)

	c.JSON(http.StatusOK, product)
}

func (ps *ProductService) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if err := ps.db.Delete(&Product{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Remove from cache
	ps.Redis.Del(context.Background(), "product:"+id)

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func (ps *ProductService) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "product-service",
		"time":    time.Now(),
	})
}
