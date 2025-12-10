# Personal Site

This is my personal website! I have a blog, projects portfolio, quotes collection, and hopefully more features in the future. I'm trying out Go, htmx, and templ for this project. Pretty cool so far. Let me know if you have any improvements for me! You can see this in action at [kgs.dev](https://kgs.dev).

## Tech Stack

- **Go** - Backend
- **Gin** - HTTP router
- **templ** - Type-safe HTML templates
- **htmx** - Dynamic HTML without JavaScript frameworks
- **PostgreSQL** - Database
- **Railway** - Hosting

## Project Structure

```
.
├── cmd/
│   ├── server/         # Main application entry point
│   └── seed/           # CLI tool to create initial admin user
├── internal/
│   ├── app/            # Application setup, dependency injection
│   ├── config/         # Environment configuration
│   ├── controllers/    # HTTP handlers
│   ├── database/       # Database connection and migrations
│   ├── middleware/     # Auth, logging, rate limiting, security headers
│   ├── models/         # Data structures
│   └── services/       # Business logic
├── migrations/         # SQL migration files
├── static/             # CSS, JS, images
└── templates/
    ├── layouts/        # Base HTML layout
    └── pages/          # Page templates (home, blog, admin, etc.)
```

## Features

- **Blog** - Markdown/HTML posts with publish/draft status
- **Projects** - Portfolio with tags, GitHub/demo links
- **Quotes** - Collection of quotes with attribution
- **Admin Panel** - Manage all content
- **User System** - Invite-based registration, session auth
- **Structured Logging** - Request tracing with correlation IDs

## Local Development

### Prerequisites

- Go 1.21+
- PostgreSQL
- [templ](https://templ.guide/) CLI

### Setup

1. Clone the repo
   ```bash
   git clone https://github.com/yourusername/personal-site.git
   cd personal-site
   ```

2. Create a PostgreSQL database
   ```bash
   createdb personal_site
   ```

3. Copy environment variables
   ```bash
   cp .env.example .env
   # Edit .env with your database URL
   ```

4. Generate templ files
   ```bash
   templ generate
   ```

5. Run the server
   ```bash
   go run ./cmd/server
   ```

6. Seed an admin user
   ```bash
   go run ./cmd/seed you@example.com "Your Name" "your-password"
   ```

7. Visit http://localhost:3000

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://dev:dev@localhost:5432/personal_site?sslmode=disable` |
| `PORT` | Server port | `3000` |
| `ENVIRONMENT` | `development` or `production` | `development` |
| `SECURE_COOKIES` | Use secure cookies (HTTPS only) | `false` |
| `BASE_URL` | Public URL for invite links | `http://localhost:3000` |
| `SESSION_DURATION_HOURS` | Session lifetime | `168` (1 week) |

## Deployment

The app is configured for Railway deployment:

1. Push to GitHub
2. Create Railway project from repo
3. Add PostgreSQL database
4. Set environment variables
5. Deploy

See `Dockerfile` for the container build.

## Security

- Session-based authentication with secure cookies
- bcrypt password hashing
- Rate limiting on login endpoint
- CSP and security headers
- CSRF protection via SameSite=Lax cookies
- Request ID tracing for debugging

## License

MIT
