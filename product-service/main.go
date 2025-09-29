package main

import (
	"log"
	middlewares "product-service/Middlewares"
	service "product-service/Service"

	"github.com/gin-gonic/gin"
)

func main() {
	ps := service.NewProductService()

	r := gin.Default()

	// CORS middleware (same as before)
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Rate limiting middleware
	r.Use(middlewares.RateLimitMiddleware(ps))

	// Public routes
	r.GET("/health", ps.HealthCheck)
	r.GET("/api/products", ps.GetAllProducts) // Allow public read access
	r.GET("/api/products/:id", ps.GetProduct)

	// Protected routes
	protected := r.Group("/api/products")
	protected.Use(middlewares.JwtMiddleware())
	{
		protected.POST("", ps.CreateProduct)
		protected.PUT("/:id", ps.UpdateProduct)
		protected.DELETE("/:id", ps.DeleteProduct)
	}

	log.Println("Product service starting on port 8082")
	r.Run(":8082")
}
