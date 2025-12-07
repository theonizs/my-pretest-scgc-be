package router

import (
	"my-go-app/internal/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewRouter ทำหน้าที่สร้าง Gin Engine และกำหนด Route ทั้งหมด
func NewRouter(heatHandler *handler.HeatHandler) *gin.Engine {
	r := gin.Default()

	// 1. Setup CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	api := r.Group("/api")
	{
		heat := api.Group("/heat-data")
		{
			heat.GET("", heatHandler.GetHeatData)
			heat.POST("", heatHandler.CreateHeatData)
		}
	}

	return r
}