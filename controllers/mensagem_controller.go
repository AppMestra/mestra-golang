package controllers

import (
	"net/http"
	"time"

	"github.com/AppMestra/mestra-golang/models"
	"github.com/AppMestra/mestra-golang/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MensagemController struct {
	db          *gorm.DB
	processador *services.ProcessadorMensagem
}

func NewMensagemController(db *gorm.DB) *MensagemController {
	return &MensagemController{
		db:          db,
		processador: services.NewProcessadorMensagem(),
	}
}

func (c *MensagemController) ProcessarMensagem(ctx *gin.Context) {
	var request models.RequestProcessar

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"erro": "Dados inválidos: " + err.Error(),
		})
		return
	}

	mensagemWhats := models.MensagemWhatsApp{
		Autor:    request.Autor,
		Conteudo: request.Mensagem,
		Data:     time.Now(),
	}

	if request.Data != "" {
		if dataParseada, err := time.Parse("02/01/2006", request.Data); err == nil {
			mensagemWhats.Data = dataParseada
		}
	}

	if err := c.db.Create(&mensagemWhats).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"erro": "Erro ao salvar mensagem: " + err.Error(),
		})
		return
	}

	entradas, err := c.processador.ProcessarMensagem(
		request.Mensagem,
		request.Autor,
		request.Data,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"erro": "Erro ao processar mensagem: " + err.Error(),
		})
		return
	}

	for _, entradaJSON := range entradas {
		dataParseada, _ := time.Parse("2006-01-02", entradaJSON.Data)
		entrada := models.Entrada{
			Valor:     entradaJSON.Valor,
			Descricao: entradaJSON.Descricao,
			Data:      dataParseada,
		}
		c.db.Create(&entrada)
	}

	ctx.JSON(http.StatusOK, entradas)
}

func (c *MensagemController) ListarMensagens(ctx *gin.Context) {
	var mensagens []models.MensagemWhatsApp

	if err := c.db.Order("created_at desc").Find(&mensagens).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"erro": "Erro ao buscar mensagens: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, mensagens)
}

func (c *MensagemController) ListarEntradas(ctx *gin.Context) {
	var entradas []models.Entrada

	if err := c.db.Order("data desc").Find(&entradas).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"erro": "Erro ao buscar entradas: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, entradas)
}
