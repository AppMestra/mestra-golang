package main

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/AppMestra/mestra-golang/config"
	"github.com/AppMestra/mestra-golang/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func startVenomBot() {
	time.Sleep(3 * time.Second)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "cd", "&&", "node", "whatsapp-bot.js")
	} else {
		cmd = exec.Command("bash", "-c", "cd && node whatsapp-bot.js")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("🤖 Iniciando VenomBot %s")

	if err := cmd.Run(); err != nil {
		log.Printf("❌ Erro ao executar VenomBot: %v", err)
	}
}

func main() {
	if err := godotenv.Load(".env.local"); err != nil {
		if err := godotenv.Load(); err != nil {
			log.Println("Arquivo .env não encontrado, usando variáveis do sistema")
		}
	}

	config.LoadConfig()

	if config.IsProductionMode() {
		gin.SetMode(gin.ReleaseMode)

		config.InitDB()

		r := gin.Default()

		r.Use(func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}

			c.Next()
		})

		routes.ConfigurarRotas(r, config.GetDB())

		go startVenomBot()

		port := config.GetPort()
		log.Printf("🚀 API Golang MESTRA iniciando na porta %s...", port)
		log.Printf("📱 Endpoint principal: POST http://localhost%s/api/processar", port)
		r.Run(port)
	}
}
