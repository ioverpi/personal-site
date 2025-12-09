package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ioverpi/personal-site/internal/config"
	"github.com/ioverpi/personal-site/internal/middleware"
	"github.com/ioverpi/personal-site/internal/services"
	"github.com/ioverpi/personal-site/templates/pages/admin"
)

// Note: CSRF tokens removed - using SameSite=Lax cookies for CSRF protection instead

type AdminController struct {
	content  *services.AdminService
	blog     *services.BlogService
	projects *services.ProjectsService
	quotes   *services.QuotesService
	auth     *services.AuthService
	users    *services.UserService
	config   *config.Config
}

func NewAdminController(
	contentService *services.AdminService,
	blogService *services.BlogService,
	projectsService *services.ProjectsService,
	quotesService *services.QuotesService,
	authService *services.AuthService,
	userService *services.UserService,
	cfg *config.Config,
) *AdminController {
	return &AdminController{
		content:  contentService,
		blog:     blogService,
		projects: projectsService,
		quotes:   quotesService,
		auth:     authService,
		users:    userService,
		config:   cfg,
	}
}

// Auth

func (c *AdminController) LoginPage(ctx *gin.Context) {
	admin.Login("").Render(ctx.Request.Context(), ctx.Writer)
}

func (c *AdminController) Login(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")

	user, err := c.auth.Authenticate(email, password)
	if err != nil {
		admin.Login("Invalid email or password").Render(ctx.Request.Context(), ctx.Writer)
		return
	}

	// Create session
	duration := time.Duration(c.config.SessionDurationHours) * time.Hour
	session, err := c.auth.CreateSession(user.ID, duration)
	if err != nil {
		admin.Login("Failed to create session").Render(ctx.Request.Context(), ctx.Writer)
		return
	}

	// Set cookie
	maxAge := c.config.SessionDurationHours * 3600
	middleware.SetSessionCookie(ctx, session.Token, maxAge, c.config.SecureCookies)

	ctx.Redirect(http.StatusFound, "/admin")
}

func (c *AdminController) Logout(ctx *gin.Context) {
	// Delete session from database
	if token, err := ctx.Cookie(middleware.SessionCookieName); err == nil {
		c.auth.DeleteSession(token)
	}

	// Clear cookie
	middleware.SetSessionCookie(ctx, "", -1, c.config.SecureCookies)
	ctx.Redirect(http.StatusFound, "/admin/login")
}

// Dashboard

func (c *AdminController) Dashboard(ctx *gin.Context) {
	user := middleware.GetUser(ctx)
	posts, _ := c.blog.GetAllPosts()
	projects, _ := c.projects.GetAllProjects()
	quotes, _ := c.quotes.GetAllQuotes()

	admin.Dashboard(user, posts, projects, quotes).Render(ctx.Request.Context(), ctx.Writer)
}

// Users

func (c *AdminController) UsersList(ctx *gin.Context) {
	user := middleware.GetUser(ctx)
	users, _ := c.users.GetAllUsers()
	invites, _ := c.auth.GetPendingInvites()

	admin.UsersList(user, users, invites).Render(ctx.Request.Context(), ctx.Writer)
}

func (c *AdminController) NewInvite(ctx *gin.Context) {
	user := middleware.GetUser(ctx)
	admin.InviteForm(user, "").Render(ctx.Request.Context(), ctx.Writer)
}

func (c *AdminController) CreateInvite(ctx *gin.Context) {
	user := middleware.GetUser(ctx)
	email := ctx.PostForm("email")

	// Invite valid for 7 days
	invite, err := c.auth.CreateInvite(email, user.ID, 7*24*time.Hour)
	if err != nil {
		admin.InviteForm(user, "Failed to create invite").Render(ctx.Request.Context(), ctx.Writer)
		return
	}

	// Show the invite URL so admin can share it
	inviteURL := c.config.BaseURL + "/register?token=" + invite.Token
	admin.InviteSuccess(user, invite, inviteURL).Render(ctx.Request.Context(), ctx.Writer)
}

func (c *AdminController) DeleteInvite(ctx *gin.Context) {
	id := getIDParam(ctx, "id")
	c.auth.DeleteInvite(id)
	ctx.Redirect(http.StatusFound, "/admin/users")
}

// Registration (via invite)

func (c *AdminController) RegisterPage(ctx *gin.Context) {
	token := ctx.Query("token")
	invite, err := c.auth.GetInvite(token)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Invalid or expired invite")
		return
	}

	admin.Register(invite, "").Render(ctx.Request.Context(), ctx.Writer)
}

func (c *AdminController) Register(ctx *gin.Context) {
	token := ctx.PostForm("token")
	name := ctx.PostForm("name")
	password := ctx.PostForm("password")

	invite, err := c.auth.GetInvite(token)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Invalid or expired invite")
		return
	}

	// Create user
	user, err := c.users.CreateUser(services.CreateUserInput{
		Email: invite.Email,
		Name:  name,
		Role:  "user",
	})
	if err != nil {
		admin.Register(invite, "Failed to create account").Render(ctx.Request.Context(), ctx.Writer)
		return
	}

	// Create login
	_, err = c.auth.CreatePasswordLogin(user.ID, invite.Email, password)
	if err != nil {
		admin.Register(invite, "Failed to set password").Render(ctx.Request.Context(), ctx.Writer)
		return
	}

	// Mark invite as used
	c.auth.UseInvite(token)

	// Create session and log in
	duration := time.Duration(c.config.SessionDurationHours) * time.Hour
	session, _ := c.auth.CreateSession(user.ID, duration)
	maxAge := c.config.SessionDurationHours * 3600
	middleware.SetSessionCookie(ctx, session.Token, maxAge, c.config.SecureCookies)

	ctx.Redirect(http.StatusFound, "/admin")
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

	_, err := c.content.CreatePost(input)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/admin")
}

func (c *AdminController) EditPost(ctx *gin.Context) {
	id := getIDParam(ctx, "id")
	post, err := c.blog.GetPostByID(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	admin.PostEditor(post).Render(ctx.Request.Context(), ctx.Writer)
}

func (c *AdminController) UpdatePost(ctx *gin.Context) {
	id := getIDParam(ctx, "id")
	input := services.UpdatePostInput{
		Title:   ctx.PostForm("title"),
		Slug:    ctx.PostForm("slug"),
		Content: ctx.PostForm("content"),
		Publish: ctx.PostForm("publish") == "on",
	}

	if _, err := c.content.UpdatePost(id, input); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/admin")
}

func (c *AdminController) DeletePost(ctx *gin.Context) {
	id := getIDParam(ctx, "id")
	c.content.DeletePost(id)
	ctx.Redirect(http.StatusFound, "/admin")
}

// Projects

func (c *AdminController) NewProject(ctx *gin.Context) {
	admin.ProjectEditor(nil).Render(ctx.Request.Context(), ctx.Writer)
}

func (c *AdminController) CreateProject(ctx *gin.Context) {
	displayOrder, _ := strconv.Atoi(ctx.PostForm("display_order")) // defaults to 0 if invalid
	tags := parseTags(ctx.PostForm("tags"))

	input := services.CreateProjectInput{
		Name:         ctx.PostForm("name"),
		Description:  ctx.PostForm("description"),
		Tags:         tags,
		GithubURL:    ctx.PostForm("github_url"),
		DemoURL:      ctx.PostForm("demo_url"),
		DisplayOrder: displayOrder,
	}

	_, err := c.content.CreateProject(input)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/admin")
}

func (c *AdminController) EditProject(ctx *gin.Context) {
	id := getIDParam(ctx, "id")
	project, err := c.projects.GetProjectByID(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	admin.ProjectEditor(project).Render(ctx.Request.Context(), ctx.Writer)
}

func (c *AdminController) UpdateProject(ctx *gin.Context) {
	id := getIDParam(ctx, "id")
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

	if _, err := c.content.UpdateProject(id, input); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/admin")
}

func (c *AdminController) DeleteProject(ctx *gin.Context) {
	id := getIDParam(ctx, "id")
	c.content.DeleteProject(id)
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

	_, err := c.content.CreateQuote(input)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/admin")
}

func (c *AdminController) EditQuote(ctx *gin.Context) {
	id := getIDParam(ctx, "id")
	quote, err := c.quotes.GetQuoteByID(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	admin.QuoteEditor(quote).Render(ctx.Request.Context(), ctx.Writer)
}

func (c *AdminController) UpdateQuote(ctx *gin.Context) {
	id := getIDParam(ctx, "id")
	input := services.UpdateQuoteInput{
		Content: ctx.PostForm("content"),
		Author:  ctx.PostForm("author"),
		IsOwn:   ctx.PostForm("is_own") == "on",
	}

	if _, err := c.content.UpdateQuote(id, input); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusFound, "/admin")
}

func (c *AdminController) DeleteQuote(ctx *gin.Context) {
	id := getIDParam(ctx, "id")
	c.content.DeleteQuote(id)
	ctx.Redirect(http.StatusFound, "/admin")
}

// Helpers

// getIDParam extracts and validates an integer ID from URL params.
// Panics on invalid ID (caught by Gin's recovery middleware).
func getIDParam(ctx *gin.Context, name string) int {
	id, err := strconv.Atoi(ctx.Param(name))
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		panic(http.StatusBadRequest)
	}
	return id
}

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
