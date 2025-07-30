package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	From      string         `json:"from" gorm:"column:from_number;not null"`
	To        string         `json:"to" gorm:"column:to_number;not null"`
	Body      string         `json:"body" gorm:"not null"`
	Type      string         `json:"type" gorm:"default:text"` // text, image, audio, document
	MediaURL  string         `json:"media_url,omitempty" gorm:"column:media_url"`
	DataHora  time.Time      `json:"data_hora" gorm:"default:CURRENT_TIMESTAMP"`
	Entrada   bool           `json:"entrada" gorm:"default:true"` // true = recebida, false = enviada
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Message) TableName() string {
	return "messages"
}

type WhatsAppWebhook struct {
	From     string `json:"from" binding:"required"`
	To       string `json:"to" binding:"required"`
	Body     string `json:"body" binding:"required"`
	Type     string `json:"type"`
	MediaURL string `json:"media_url,omitempty"`
}

type SendMessageRequest struct {
	To      string `json:"to" binding:"required"`
	Message string `json:"message" binding:"required"`
	Type    string `json:"type,omitempty"`
}
