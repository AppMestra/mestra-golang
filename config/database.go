package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error

	dbConfig := AppConfig.Database

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name, dbConfig.SSLMode)

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		dsn = dbURL
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	if IsProductionMode() {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatal("Erro ao conectar com o banco de dados:", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Erro ao obter instância SQL:", err)
	}

	if err = sqlDB.Ping(); err != nil {
		log.Fatal("Erro ao fazer ping no banco de dados:", err)
	}

	fmt.Println("Conectado ao banco de dados PostgreSQL com GORM!")

	autoMigrate()
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func autoMigrate() {
	type Message struct {
		ID        uint           `json:"id" gorm:"primaryKey"`
		From      string         `json:"from" gorm:"column:from_number;not null"`
		To        string         `json:"to" gorm:"column:to_number;not null"`
		Body      string         `json:"body" gorm:"not null"`
		Type      string         `json:"type" gorm:"default:text"`
		MediaURL  string         `json:"media_url,omitempty" gorm:"column:media_url"`
		DataHora  time.Time      `json:"data_hora" gorm:"default:CURRENT_TIMESTAMP"`
		Entrada   bool           `json:"entrada" gorm:"default:true"`
		CreatedAt time.Time      `json:"created_at"`
		UpdatedAt time.Time      `json:"updated_at"`
		DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	}

	err := DB.AutoMigrate(&Message{})
	if err != nil {
		log.Fatal("Erro ao executar auto-migração:", err)
	}

	fmt.Println("Auto-migração executada com sucesso!")
}
