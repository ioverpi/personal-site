package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ioverpi/personal-site/templates/pages"
)

type HomeController struct{}

func NewHomeController() *HomeController {
	return &HomeController{}
}

func (c *HomeController) Index(ctx *gin.Context) {
	pages.Home().Render(ctx.Request.Context(), ctx.Writer)
}
