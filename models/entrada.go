package models

import (
	"time"
)

type Entrada struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Valor     float64   `json:"valor"`
	Descricao string    `json:"descricao"`
	Data      time.Time `json:"data"`
	CreatedAt time.Time `json:"created_at"`
}

type EntradaJSON struct {
	Valor     float64 `json:"valor"`
	Descricao string  `json:"descricao"`
	Data      string  `json:"data"` // formato YYYY-MM-DD
}

type MensagemWhatsApp struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Autor     string    `json:"autor"`
	Data      time.Time `json:"data"`
	Conteudo  string    `json:"conteudo"`
	CreatedAt time.Time `json:"created_at"`
}

type RequestProcessar struct {
	Mensagem string `json:"mensagem" binding:"required"`
	Autor    string `json:"autor"`
	Data     string `json:"data"`
}
