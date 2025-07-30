package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AppMestra/mestra-golang/config"
	"github.com/AppMestra/mestra-golang/models"
)

func SendWhatsAppMessage(to, message, msgType string) error {
	// Aqui você integraria com VenomBot ou outra biblioteca WhatsApp
	// Por enquanto, vamos simular o envio e salvar no banco

	if msgType == "" {
		msgType = "text"
	}

	msg := models.Message{
		From:     "bot",
		To:       to,
		Body:     message,
		Type:     msgType,
		DataHora: time.Now(),
		Entrada:  false,
	}

	return config.DB.Create(&msg).Error
}

func ProcessIncomingMessage(webhook models.WhatsAppWebhook) error {
	msg := models.Message{
		From:     webhook.From,
		To:       webhook.To,
		Body:     webhook.Body,
		Type:     webhook.Type,
		MediaURL: webhook.MediaURL,
		DataHora: time.Now(),
		Entrada:  true,
	}

	if err := config.DB.Create(&msg).Error; err != nil {
		return err
	}

	return forwardToAPI(webhook)
}

func GetMessages(limit int) ([]models.Message, error) {
	var messages []models.Message

	err := config.DB.
		Order("data_hora DESC").
		Limit(limit).
		Find(&messages).Error

	return messages, err
}

func forwardToAPI(webhook models.WhatsAppWebhook) error {
	apiURL := config.AppConfig.API.DotNetURL

	jsonData, err := json.Marshal(webhook)
	if err != nil {
		return err
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API retornou status: %d", resp.StatusCode)
	}

	return nil
}
