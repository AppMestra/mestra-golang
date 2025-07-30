package main

import (
	"log"

	"github.com/AppMestra/mestra-golang/config"
	"github.com/AppMestra/mestra-golang/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env.local"); err != nil {
		if err := godotenv.Load(); err != nil {
			log.Println("Arquivo .env não encontrado, usando variáveis do sistema")
		}
	}

	config.LoadConfig()

	if config.IsProductionMode() {
		gin.SetMode(gin.ReleaseMode)
	}

	config.InitDB()

	r := gin.Default()

	routes.SetupRoutes(r)

	port := config.GetPort()
	log.Printf("Servidor WhatsApp Bot iniciando na porta %s...", port)
	r.Run(port)
}
