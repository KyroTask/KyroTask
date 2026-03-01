# 🦊 KyroTask

<div align="center">

A premium, full-stack productivity ecosystem and habit tracker for Telegram Mini App built with **Vue 3** and **Go**.

[![Backend](https://img.shields.io/badge/Backend-Go-00ADD8?style=flat-square&logo=go&logoColor=white)](https://github.com/KyroTask/internal)
[![Frontend](https://img.shields.io/badge/Frontend-Vue%203-4FC08D?style=flat-square&logo=vue.js&logoColor=white)](https://github.com/KyroTask/mini-app)
[![Platform](https://img.shields.io/badge/Platform-Telegram%20Mini%20App-2CA5E0?style=flat-square&logo=telegram&logoColor=white)](https://telegram.org)
[![License](https://img.shields.io/badge/License-MIT-yellow?style=flat-square)](LICENSE)

</div>

---

## 📦 Repositories

This is the **main monorepo**. Core packages live in their own repos linked as git submodules:

| Repo | Description | Stack |
|------|-------------|-------|
| [**KyroTask**](https://github.com/KyroTask/KyroTask) | Root config, entry point, scripts | — |
| [**mini-app**](https://github.com/KyroTask/mini-app) | Telegram Mini App frontend | Vue 3 · Vite · Tailwind · Pinia |
| [**internal**](https://github.com/KyroTask/internal) | REST API + WebSocket backend | Go · Gin · GORM · PostgreSQL |

---

## ✨ Features

- ✅ **Task Management** — Create, edit, complete, and organize tasks with subtasks and comments
- 📁 **Projects & Milestones** — Group tasks with automated progress and milestone tracking
- 🎯 **Goal Tracking** — Set and monitor long-term goals
- 🔥 **Habit Tracker** — Build streaks and track daily habits
- ⏱️ **Pomodoro Focus Timer** — Deep work sessions with automatic cycle tracking and XP leveling
- 📊 **Advanced Analytics** — Visual dashboard for focus time, task completion, and habits
- 📅 **Calendar View** — Visualize tasks and deadlines by date
- 🔔 **Telegram Notifications** — Automated reminders and alerts via bot
- 🔐 **Secure Auth** — Telegram WebApp + Google (Firebase) authentication
- 📱 **Mobile First** — Premium UI optimized for Telegram mobile and desktop web

---

## 🛠 Tech Stack

### Frontend ([mini-app](https://github.com/KyroTask/mini-app))

| Tool | Purpose |
|------|---------|
| Vue 3 + Composition API | UI framework |
| Vite 7 | Build tool & dev server |
| Tailwind CSS 3 | Utility-first styling |
| Pinia | State management |
| Vue Router 4 | Client-side routing |
| Axios | HTTP client |
| @telegram-apps/sdk | Telegram WebApp integration |
| Firebase 12 | Google authentication |

### Backend ([internal](https://github.com/KyroTask/internal))

| Tool | Purpose |
|------|---------|
| Go 1.25 | Backend language |
| Gin | HTTP web framework |
| GORM | ORM |
| PostgreSQL / SQLite | Database |
| JWT (golang-jwt v5) | Auth tokens |
| Firebase Admin SDK | Google auth verification |
| Telegram Bot API | Bot notifications & webhooks |
| Gorilla WebSocket | Real-time Pomodoro sync |

---

## 📂 Project Structure

```text
KyroTask/
├── cmd/
│   └── server/
│       └── main.go              # Go application entry point
├── internal/                    # → submodule: KyroTask/internal
│   ├── config/                  # App configuration
│   ├── database/                # DB connection & migrations
│   ├── handlers/                # HTTP request handlers
│   ├── middleware/              # Auth, rate limit, security
│   ├── models/                  # GORM database models
│   └── services/                # Business logic & scheduler
├── mini-app/                    # → submodule: KyroTask/mini-app
│   ├── src/
│   │   ├── pages/               # Route-level page components
│   │   ├── components/          # Reusable UI components
│   │   ├── stores/              # Pinia state stores
│   │   ├── router/              # Vue Router config
│   │   └── services/            # API client (Axios)
│   └── package.json
├── .env.example                 # Environment template
├── .air.toml                    # Air hot-reload config (Go)
├── go.mod                       # Go module definition
├── package.json                 # Root npm scripts
└── nixpacks.toml                # Cloud deployment config
```

---

## 🚀 Getting Started

### Prerequisites

- **Go** 1.21+
- **Node.js** 18+
- **air** (Go hot-reload): `go install github.com/air-verse/air@latest`

### Clone (with submodules)

```bash
git clone --recurse-submodules https://github.com/KyroTask/KyroTask.git
cd KyroTask
```

### Install Dependencies

```bash
npm install          # Root runner (concurrently)
go mod download      # Go dependencies
cd mini-app && npm install && cd ..   # Frontend
```

### Configure Environment

```bash
cp .env.example .env
```

Edit `.env` and set at minimum:

```env
JWT_SECRET=your-secure-random-string
TELEGRAM_BOT_TOKEN=your-bot-token
DB_DRIVER=sqlite
DB_DSN=./data/dev.db
```

### Start Development

```bash
npm run dev
```

| Service | URL |
|---------|-----|
| Backend (Go + air) | `http://localhost:3001` |
| Frontend (Vite) | `http://localhost:5173` |

Frontend proxies all `/api` requests to the backend automatically.

### Individual Commands

```bash
npm run dev:server     # Backend only
npm run dev:client     # Frontend only
npm run build          # Build frontend
npm run build:all      # Build frontend + Go binary
```

---

## 🔐 Telegram Bot Setup

1. Create a bot via [@BotFather](https://t.me/botfather) using `/newbot`
2. Save the bot token → set `TELEGRAM_BOT_TOKEN` in `.env`
3. Create a Mini App (`/newapp`) and point the URL to your deployment
4. Register the webhook:

```bash
curl -X POST https://api.telegram.org/bot<TOKEN>/setWebhook \
  -d "url=https://your-domain.com/api/v1/telegram/webhook"
```

---

## 📝 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `3001` |
| `ENV` | `development` or `production` | `development` |
| `DB_DRIVER` | `sqlite` or `postgres` | `sqlite` |
| `DB_DSN` | Database path / connection string | `./data/dev.db` |
| `JWT_SECRET` | Secret key for JWT signing | **required** |
| `JWT_EXPIRY` | JWT token lifetime | `168h` |
| `TELEGRAM_BOT_TOKEN` | Telegram bot token | **required** |
| `ALLOWED_ORIGINS` | CORS allowed origins | `http://localhost:5173` |

---

## 🏗 Production Deployment

```bash
npm run build:all
# Produces a single Go binary at bin/server with embedded Vue SPA
```

Supported deployment platforms (via `nixpacks.toml`):

- **Railway** / **Render** / **Fly.io** — push and deploy directly
- **DigitalOcean** App Platform
- **AWS** ECS / Elastic Beanstalk
- Any platform supporting Go binaries

---

## 📋 Roadmap

- [x] Authentication (Telegram WebApp + Google)
- [x] Database models & migrations
- [x] Projects & Tasks CRUD (subtasks, comments)
- [x] Goals tracking
- [x] Habits tracker with streaks
- [x] Calendar view
- [x] Activity feed
- [x] Telegram notifications & reminders
- [x] Milestones
- [x] Analytics dashboard
- [x] Pomodoro focus timer with XP leveling
- [x] WebSocket real-time sync
- [ ] E2E tests
- [ ] Production deployment

---

## 🤝 Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## 📄 License

MIT © [KyroTask](https://github.com/KyroTask)
