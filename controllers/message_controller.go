package controllers

import (
	"net/http"
	"strconv"

	"github.com/AppMestra/mestra-golang/models"
	"github.com/AppMestra/mestra-golang/services"
	"github.com/gin-gonic/gin"
)

func SendMessage(c *gin.Context) {
	var req models.SendMessageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.SendWhatsAppMessage(req.To, req.Message, req.Type); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao enviar mensagem"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mensagem enviada com sucesso"})
}

func ReceiveWebhook(c *gin.Context) {
	var webhook models.WhatsAppWebhook

	if err := c.ShouldBindJSON(&webhook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.ProcessIncomingMessage(webhook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar mensagem"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func GetMessages(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	messages, err := services.GetMessages(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar mensagens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}
