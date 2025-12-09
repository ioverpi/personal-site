# Personal Website Implementation Plan

## Overview
Minimalist developer portfolio with blog, projects, and quotes pages. Built with Go + templ + htmx, backed by Postgres, deployed via Docker.

## Tech Stack
- **Backend**: Go + Gin
- **Templates**: templ
- **Interactivity**: htmx + hx-boost for smooth navigation
- **Database**: PostgreSQL
- **Styling**: Plain CSS (custom properties for theming)
- **Local Dev**: docker-compose (Postgres)
- **Production**: Railway (app + Postgres plugin)

---

## Project Structure

```
personal_site/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── app/                     # Dependency container
│   │   └── app.go               # App struct, wires everything together
│   ├── config/                  # Configuration
│   │   └── config.go            # Load from env vars
│   ├── controllers/             # HTTP layer (params, response, status)
│   │   ├── blog.go
│   │   ├── projects.go
│   │   ├── quotes.go
│   │   └── admin.go
│   ├── services/                # Business logic layer
│   │   ├── blog.go
│   │   ├── projects.go
│   │   └── quotes.go
│   ├── models/                  # Database models
│   │   ├── post.go
│   │   ├── project.go
│   │   └── quote.go
│   ├── database/                # DB connection
│   │   └── db.go
│   ├── adapters/                # Future: external service implementations
│   │   └── .gitkeep             # (empty for now, add providers later)
│   └── middleware/              # Auth, logging, etc.
│       └── middleware.go
├── templates/                   # templ files
│   ├── layouts/
│   │   └── base.templ           # Shared layout
│   ├── pages/
│   │   ├── home.templ
│   │   ├── blog.templ
│   │   ├── post.templ
│   │   ├── projects.templ
│   │   ├── quotes.templ
│   │   └── admin/
│   │       ├── dashboard.templ
│   │       ├── post_editor.templ
│   │       ├── project_editor.templ
│   │       └── quote_editor.templ
│   └── components/
│       ├── nav.templ
│       ├── footer.templ
│       ├── post_card.templ
│       ├── project_card.templ
│       └── quote_card.templ
├── static/
│   ├── css/
│   │   └── style.css
│   └── js/
│       └── theme.js             # Dark/light toggle
├── migrations/                  # SQL migrations
│   ├── 001_create_posts.sql
│   ├── 002_create_projects.sql
│   └── 003_create_quotes.sql
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum
```

---

## Database Schema

### posts
| Column       | Type         | Notes                           |
|--------------|--------------|-------------------------------- |
| id           | SERIAL       | Primary key                     |
| title        | VARCHAR(255) |                                 |
| slug         | VARCHAR(255) | URL-friendly, unique            |
| content      | TEXT         | Markdown content                |
| published_at | TIMESTAMP    | NULL = draft, set = published   |
| created_at   | TIMESTAMP    |                                 |
| updated_at   | TIMESTAMP    |                                 |

### projects
| Column       | Type         | Notes                    |
|--------------|--------------|--------------------------|
| id           | SERIAL       | Primary key              |
| name         | VARCHAR(255) |                          |
| description  | TEXT         |                          |
| tags         | TEXT[]       | Array of tech tags       |
| github_url   | VARCHAR(255) | Nullable, GitHub repo    |
| demo_url     | VARCHAR(255) | Nullable, live demo      |
| display_order| INT          | For manual sorting       |
| created_at   | TIMESTAMP    |                          |

### quotes
| Column       | Type         | Notes                    |
|--------------|--------------|--------------------------|
| id           | SERIAL       | Primary key              |
| content      | TEXT         | The quote itself         |
| author       | VARCHAR(255) |                          |
| is_own       | BOOLEAN      | Your quote vs collected  |
| created_at   | TIMESTAMP    |                          |

---

## Implementation Steps

### Phase 1: Foundation
1. Initialize Go module
2. Set up project directory structure
3. Create docker-compose.yml with Postgres
4. Set up database connection and migrations
5. Create base templ layout with nav

### Phase 2: Public Pages
6. Home page (intro + nav links)
7. Blog list page (fetch published posts)
8. Individual blog post page (render markdown)
9. Projects page (list with tags)
10. Quotes page (list with own/collected distinction)

### Phase 3: Styling & Interactivity
11. CSS with custom properties for theming
12. Dark/light mode toggle (localStorage + JS)
13. htmx hx-boost for instant page transitions
14. Random quote button on quotes page

### Phase 4: Admin UI
15. Simple auth (environment-based password or session)
16. Admin dashboard
17. Blog post editor (create/edit/delete, markdown preview)
18. Project editor (create/edit/delete, tag management)
19. Quote editor (create/edit/delete)

### Phase 5: Deployment
20. Dockerfile for Go app
21. docker-compose for full stack (app + postgres)
22. Environment configuration

---

## Key Files to Create

| Priority | File | Purpose |
|----------|------|---------|
| 1 | `go.mod` | Module initialization |
| 1 | `docker-compose.yml` | Local dev environment |
| 1 | `cmd/server/main.go` | App entry point, route setup |
| 2 | `internal/database/db.go` | Postgres connection |
| 2 | `migrations/*.sql` | Database schema |
| 2 | `internal/models/*.go` | Data structures |
| 3 | `internal/services/*.go` | Business logic |
| 3 | `internal/controllers/*.go` | HTTP handlers |
| 3 | `templates/layouts/base.templ` | Shared HTML structure |
| 3 | `templates/pages/*.templ` | Page templates |
| 4 | `static/css/style.css` | Styling |
| 5 | `Dockerfile` | Production container (Railway) |

---

## Notes

- Using `hx-boost` on nav links gives SPA-like transitions with minimal JS
- Markdown rendering via `goldmark` library
- Admin auth can start simple (env var password) and evolve later
- CSS custom properties make dark/light theming straightforward

## Future: Adding External Dependencies

When you need to add an external service (LLM, S3, email, etc.):

1. Create interface in `internal/<service>/` (e.g., `internal/llm/llm.go`)
2. Create adapter in `internal/adapters/<provider>/` (e.g., `internal/adapters/claude/`)
3. Add interface field to `App` struct
4. Wire up in `app.New()` based on config

Example structure when adding LLM + storage:
```
internal/
├── llm/
│   └── llm.go              # type LLM interface { Complete(...) }
├── storage/
│   └── storage.go          # type Storage interface { Upload(...) }
├── adapters/
│   ├── claude/
│   │   └── claude.go       # implements llm.LLM
│   ├── openai/
│   │   └── openai.go       # implements llm.LLM
│   ├── s3/
│   │   └── s3.go           # implements storage.Storage
│   └── local/
│       └── local.go        # implements storage.Storage
```

This pattern keeps dependencies swappable and testable.
