package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ioverpi/personal-site/internal/services"
	"github.com/ioverpi/personal-site/templates/pages"
)

type ProjectsController struct {
	projects *services.ProjectsService
}

func NewProjectsController(projects *services.ProjectsService) *ProjectsController {
	return &ProjectsController{projects: projects}
}

func (c *ProjectsController) List(ctx *gin.Context) {
	projects, err := c.projects.GetAllProjects()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	pages.ProjectsList(projects).Render(ctx.Request.Context(), ctx.Writer)
}

func (c *ProjectsController) LastGameOf2020(ctx *gin.Context) {
	pages.LastGameOf2020().Render(ctx.Request.Context(), ctx.Writer)
}
