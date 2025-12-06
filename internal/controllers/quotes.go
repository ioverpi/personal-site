package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ioverpi/personal-site/internal/services"
	"github.com/ioverpi/personal-site/templates/pages"
)

type QuotesController struct {
	quotes *services.QuotesService
}

func NewQuotesController(quotes *services.QuotesService) *QuotesController {
	return &QuotesController{quotes: quotes}
}

func (c *QuotesController) List(ctx *gin.Context) {
	quotes, err := c.quotes.GetAllQuotes()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	pages.QuotesList(quotes).Render(ctx.Request.Context(), ctx.Writer)
}

func (c *QuotesController) Random(ctx *gin.Context) {
	quote, err := c.quotes.GetRandomQuote()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	pages.QuoteCard(quote).Render(ctx.Request.Context(), ctx.Writer)
}
