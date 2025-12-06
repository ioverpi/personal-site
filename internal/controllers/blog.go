package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ioverpi/personal-site/internal/services"
	"github.com/ioverpi/personal-site/templates/pages"
)

type BlogController struct {
	blog *services.BlogService
}

func NewBlogController(blog *services.BlogService) *BlogController {
	return &BlogController{blog: blog}
}

func (c *BlogController) List(ctx *gin.Context) {
	posts, err := c.blog.GetPublishedPosts()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	pages.BlogList(posts).Render(ctx.Request.Context(), ctx.Writer)
}

func (c *BlogController) Show(ctx *gin.Context) {
	slug := ctx.Param("slug")

	post, err := c.blog.GetPostBySlug(slug)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	// Only show published posts to public
	if post.PublishedAt == nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	pages.BlogPost(post).Render(ctx.Request.Context(), ctx.Writer)
}
