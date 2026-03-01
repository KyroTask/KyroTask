# 🦊 KyroTask

A premium, full-stack productivity ecosystem and habit tracker for Telegram Mini App built with **Vue 3** (Composition API) and **Go**.

## 🚀 Features

- ✅ **Full Task Management** - Create, edit, complete, and organize tasks
- 📁 **Projects & Goals** - Group tasks logically with automated progress and milestone tracking
- 🔥 **Habit Tracker** - Build streaks and track daily habits
- ⏱️ **Pomodoro Focus Timer** - Deep work sessions with automatic cycle tracking and leveling
- 📊 **Advanced Analytics** - Visual dashboard for tracking focus time, task completion, and habits
- 📅 **Calendar View** - Visualize tasks by due date
- 🔔 **Telegram Notifications** - Automated reminders and alerts via bot
- 🔐 **Secure Auth** - Multi-provider authentication via Telegram WebApp & Google
- 📱 **Mobile First** - Premium UI/UX optimized for both Telegram mobile app and desktop web.

## 🛠 Tech Stack

### Frontend

- **Vue 3** - Composition API
- **Vite** - Build tool and dev server
- **Tailwind CSS** - Utility-first styling
- **Vue Router** - Client-side routing
- **Pinia** - State management
- **Axios** - HTTP client
- **Telegram WebApp SDK** - Telegram integration

### Backend

- **Go** - Backend language
- **Gin** - Web framework
- **GORM** - ORM
- **PostgreSQL** / **SQLite** - Database
- **JWT** - Authentication tokens

## 📦 Project Structure

```text
telegram-task-manager/
├── cmd/
│   └── server/              # Go main entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── database/            # Database connection
│   ├── handlers/            # HTTP handlers
│   ├── middleware/          # Middleware (auth, cors)
│   ├── models/              # Database models
│   └── services/            # Business logic
├── migrations/              # Database migrations
├── mini-app/                # Vue 3 frontend
│   ├── src/
│   │   ├── pages/           # Page components
│   │   ├── components/      # Reusable components
│   │   ├── stores/          # Pinia stores
│   │   ├── router/          # Vue Router
│   │   └── services/        # API services
│   └── package.json
├── .env                     # Environment variables
├── .air.toml                # Air config (Go auto-reload)
├── notifications_spec.md    # Telegram notifications specification
└── package.json             # Root scripts
```

## 🚀 Getting Started

### Prerequisites

- **Go 1.21+**
- **Node.js 18+**
- **air** (for Go auto-reload): `go install github.com/air-verse/air@latest`

### Installation

1. **Clone the repository**

```bash
git clone <repository-url>
cd telegram-task-manager
```

1. **Install dependencies**

```bash
# Install root dependencies (concurrently)
npm install

# Install Go dependencies
go mod download

# Install frontend dependencies
cd mini-app && npm install && cd ..
```

1. **Configure environment**

```bash
# Copy example env file
cp .env.example .env

# Edit .env and set:
# - JWT_SECRET (generate a secure random string)
# - TELEGRAM_BOT_TOKEN (get from @BotFather)
```

### Development

**Single command to run both backend and frontend:**

```bash
npm run dev
```

This will start:

- **Backend** (Go with air auto-reload) on `http://localhost:3001`
- **Frontend** (Vite dev server) on `http://localhost:5173`

The frontend will proxy API requests to the backend automatically.

### Individual Commands

```bash
# Run backend only
npm run dev:server

# Run frontend only
npm run dev:client

# Build frontend
npm run build

# Build everything (frontend + backend binary)
npm run build:all
```

## 🔐 Telegram Bot Setup

1. **Create a bot** with [@BotFather](https://t.me/botfather)
   - Use `/newbot` command
   - Save the bot token

2. **Create a Mini App**
   - Use `/newapp` command with @BotFather
   - Set the Mini App URL to your deployment URL
   - For local development, use tools like ngrok

3. **Set webhook** (for bot commands)

```bash
curl -X POST https://api.telegram.org/bot<YOUR_BOT_TOKEN>/setWebhook \
  -d "url=https://your-domain.com/api/v1/telegram/webhook"
```

## 📝 Environment Variables

| Variable             | Description                              | Default                                       |
| -------------------- | ---------------------------------------- | --------------------------------------------- |
| `PORT`               | Server port                              | `3001`                                        |
| `DB_DRIVER`          | Database driver (`sqlite` or `postgres`) | `sqlite`                                      |
| `DB_DSN`             | Database connection string               | `./data/dev.db`                               |
| `JWT_SECRET`         | Secret key for JWT signing               | _required_                                    |
| `JWT_EXPIRY`         | JWT token expiration                     | `168h` (7 days)                               |
| `TELEGRAM_BOT_TOKEN` | Telegram bot token                       | _required_                                    |
| `ALLOWED_ORIGINS`    | CORS allowed origins                     | `http://localhost:5173,http://localhost:3001` |

## 🏗 Production Deployment

### Build

```bash
npm run build:all
```

This creates a single Go binary at `bin/server` with the Vue SPA embedded.

### Docker

```dockerfile
docker build -t telegram-task-manager .
docker run -p 3001:3001 \
  -e JWT_SECRET=your-secret \
  -e TELEGRAM_BOT_TOKEN=your-token \
  -e DB_DRIVER=postgres \
  -e DB_DSN=postgres://... \
  telegram-task-manager
```

### Deploy to Cloud

The application can be deployed to:

- **DigitalOcean** App Platform
- **AWS** Elastic Beanstalk / ECS
- **Fly.io**
- **Railway**
- Any platform supporting Go applications

## 🧪 Testing

```bash
# Run Go tests
go test ./internal/...

# Run frontend tests
cd mini-app && npm run test
```

## 📖 API Documentation

### Authentication

- `POST /api/v1/auth/telegram/verify` - Verify Telegram WebApp initData, returns JWT

### Protected Endpoints

All endpoints require `Authorization: Bearer <jwt>` header.

- `GET /api/v1/dashboard` - Dashboard data
- `GET /api/v1/projects` - List projects
- `POST /api/v1/projects` - Create project
- `GET /api/v1/tasks` - List tasks
- `POST /api/v1/tasks` - Create task
- ...more endpoints to be implemented

## 📋 Roadmap

- [x] Project setup and scaffolding
- [x] Basic authentication structure
- [x] Database models
- [x] Frontend routing and pages
- [x] Telegram WebApp auth implementation
- [x] Projects CRUD
- [x] Tasks CRUD with subtasks and comments
- [x] Goals tracking
- [x] Habits tracker
- [x] Calendar view
- [x] Activity logging
- [x] Telegram Notifications & Reminders
- [x] Milestones implementation
- [x] Reports and Analytics Dashboard
- [x] Pomodoro Focus Timer
- [ ] E2E tests
- [ ] Production deployment

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

MIT License
