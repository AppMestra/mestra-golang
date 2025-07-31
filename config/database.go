package config

import (
	"fmt"
	"log"
	"os"

	"github.com/AppMestra/mestra-golang/models"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadConfig() {
	if os.Getenv("DB_HOST") == "" {
		os.Setenv("DB_HOST", "localhost")
	}
	if os.Getenv("DB_PORT") == "" {
		os.Setenv("DB_PORT", "1433")
	}
	if os.Getenv("DB_NAME") == "" {
		os.Setenv("DB_NAME", "apiGo")
	}
	if os.Getenv("DB_USER") == "" {
		os.Setenv("DB_USER", "sa")
	}
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", ":8080")
	}
}

func InitDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		username, password, host, port, database)

	var err error
	DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar com o banco de dados:", err)
	}

	err = DB.AutoMigrate(&models.MensagemWhatsApp{}, &models.Entrada{})
	if err != nil {
		log.Fatal("Erro na migração das tabelas:", err)
	}

	log.Println("Conectado ao banco de dados SQL Server com sucesso!")
}

func GetDB() *gorm.DB {
	return DB
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ":8080"
	}
	if port[0] != ':' {
		return ":" + port
	}
	return port
}

func IsProductionMode() bool {
	return os.Getenv("GIN_MODE") == "release"
}
