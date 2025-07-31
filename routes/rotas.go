package routes

import (
	"github.com/AppMestra/mestra-golang/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ConfigurarRotas(r *gin.Engine, db *gorm.DB) {
	mensagemController := controllers.NewMensagemController(db)

	api := r.Group("/api")
	{
		api.POST("/processar", mensagemController.ProcessarMensagem)
		api.GET("/mensagens", mensagemController.ListarMensagens)
		api.GET("/entradas", mensagemController.ListarEntradas)
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API Golang MESTRA funcionando!",
		})
	})
}
