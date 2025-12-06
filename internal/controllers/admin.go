package controllers

import (
	"crypto/subtle"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ioverpi/personal-site/internal/services"
	"github.com/ioverpi/personal-site/templates/pages/admin"
)

type AdminController struct {
	admin    *services.AdminService
	blog     *services.BlogService
	projects *services.ProjectsService
	quotes   *services.QuotesService
	password string
}

func NewAdminController(
	adminService *services.AdminService,
	blogService *services.BlogService,
	projectsService *services.ProjectsService,
	quotesService *services.QuotesService,
	password string,
) *AdminController {
	return &AdminController{
		admin:    adminService,
		blog:     blogService,
		projects: projectsService,
		quotes:   quotesService,
		password: password,
	}
}

// Auth

func (c *AdminController) LoginPage(ctx *gin.Context) {
	admin.Login("").Render(ctx.Request.Context(), ctx.Writer)
}

func (c *AdminController) Login(ctx *gin.Context) {
	password := ctx.PostForm("password")

	if subtle.ConstantTimeCompare([]byte(password), []byte(c.password)) == 1 {
		ctx.SetCookie("admin_session", c.password, 86400*7, "/", "", false, true)
		ctx.Redirect(http.StatusFound, "/admin")
		return
	}

	admin.Login("Invalid password").Render(ctx.Request.Context(), ctx.Writer)
}

func (c *AdminController) Logout(ctx *gin.Context) {
	ctx.SetCookie("admin_session", "", -1, "/", "", false, true)
	ctx.Redirect(http.StatusFound, "/admin/login")
}

// Dashboard

func (c *AdminController) Dashboard(ctx *gin.Context) {
	posts, _ := c.blog.GetAllPosts()
	projects, _ := c.projects.GetAllProjects()
	quotes, _ := c.quotes.GetAllQuotes()

	admin.Dashboard(posts, projects, quotes).Render(ctx.Request.Context(), ctx.Writer)
}

// Posts

func (c *AdminController) NewPost(ctx *gin.Context) {
	admin.PostEditor(nil).Render(ctx.Request.Context(), ctx.Writer)
}

func (c *AdminController) CreatePost(ctx *gin.Context) {
	input := services.CreatePostInput{
		Title:   ctx.PostForm("title"),
		Slug:    ctx.PostForm("slug"),
		Content: ctx.PostForm("content"),
		Publish: ctx.PostForm("publish") == "on",
	}

	_, err := c.admin.CreatePost(input)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/admin")
}

func (c *AdminController) EditPost(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	post, err := c.blog.GetPostByID(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	admin.PostEditor(post).Render(ctx.Request.Context(), ctx.Writer)
}

func (c *AdminController) UpdatePost(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	input := services.UpdatePostInput{
		Title:   ctx.PostForm("title"),
		Slug:    ctx.PostForm("slug"),
		Content: ctx.PostForm("content"),
		Publish: ctx.PostForm("publish") == "on",
	}

	_, err := c.admin.UpdatePost(id, input)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/admin")
}

func (c *AdminController) DeletePost(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	c.admin.DeletePost(id)
	ctx.Redirect(http.StatusFound, "/admin")
}

// Projects

func (c *AdminController) NewProject(ctx *gin.Context) {
	admin.ProjectEditor(nil).Render(ctx.Request.Context(), ctx.Writer)
}

func (c *AdminController) CreateProject(ctx *gin.Context) {
	displayOrder, _ := strconv.Atoi(ctx.PostForm("display_order"))
	tags := parseTags(ctx.PostForm("tags"))

	input := services.CreateProjectInput{
		Name:         ctx.PostForm("name"),
		Description:  ctx.PostForm("description"),
		Tags:         tags,
		GithubURL:    ctx.PostForm("github_url"),
		DemoURL:      ctx.PostForm("demo_url"),
		DisplayOrder: displayOrder,
	}

	_, err := c.admin.CreateProject(input)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/admin")
}

func (c *AdminController) EditProject(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	project, err := c.projects.GetProjectByID(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	admin.ProjectEditor(project).Render(ctx.Request.Context(), ctx.Writer)
}

func (c *AdminController) UpdateProject(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	displayOrder, _ := strconv.Atoi(ctx.PostForm("display_order"))
	tags := parseTags(ctx.PostForm("tags"))

	input := services.UpdateProjectInput{
		Name:         ctx.PostForm("name"),
		Description:  ctx.PostForm("description"),
		Tags:         tags,
		GithubURL:    ctx.PostForm("github_url"),
		DemoURL:      ctx.PostForm("demo_url"),
		DisplayOrder: displayOrder,
	}

	_, err := c.admin.UpdateProject(id, input)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/admin")
}

func (c *AdminController) DeleteProject(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	c.admin.DeleteProject(id)
	ctx.Redirect(http.StatusFound, "/admin")
}

// Quotes

func (c *AdminController) NewQuote(ctx *gin.Context) {
	admin.QuoteEditor(nil).Render(ctx.Request.Context(), ctx.Writer)
}

func (c *AdminController) CreateQuote(ctx *gin.Context) {
	input := services.CreateQuoteInput{
		Content: ctx.PostForm("content"),
		Author:  ctx.PostForm("author"),
		IsOwn:   ctx.PostForm("is_own") == "on",
	}

	_, err := c.admin.CreateQuote(input)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/admin")
}

func (c *AdminController) EditQuote(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	quote, err := c.quotes.GetQuoteByID(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	admin.QuoteEditor(quote).Render(ctx.Request.Context(), ctx.Writer)
}

func (c *AdminController) UpdateQuote(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	input := services.UpdateQuoteInput{
		Content: ctx.PostForm("content"),
		Author:  ctx.PostForm("author"),
		IsOwn:   ctx.PostForm("is_own") == "on",
	}

	_, err := c.admin.UpdateQuote(id, input)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/admin")
}

func (c *AdminController) DeleteQuote(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	c.admin.DeleteQuote(id)
	ctx.Redirect(http.StatusFound, "/admin")
}

// Helpers

func parseTags(tagsStr string) []string {
	if tagsStr == "" {
		return []string{}
	}

	parts := strings.Split(tagsStr, ",")
	tags := make([]string, 0, len(parts))
	for _, tag := range parts {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			tags = append(tags, tag)
		}
	}
	return tags
}
