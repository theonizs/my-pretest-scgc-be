package main

import (
	"my-go-app/internal/config"
	"my-go-app/internal/handler"
	"my-go-app/internal/repository"
	"my-go-app/internal/router"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Setup Infrastructure (Database)
	repo := repository.NewHeatRepository()

	// 3. Setup Handler (Controller)
	heatHandler := handler.NewHeatHandler(repo)

	// 4. Init Router (ย้าย logic ไปไว้ใน package router แล้ว)
	r := router.NewRouter(heatHandler)

	// 5. Run Server
	r.Run(cfg.ServerPort)
}