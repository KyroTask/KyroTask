# Quick Start Commands

## 🚀 Start Everything (Easiest)

```bash
./start.sh
```

This script will:

- Check and create `.env` if needed
- Install `air` if not present
- Install all dependencies
- Start both backend and frontend servers

---

## 📝 Manual Commands

### First Time Setup

```bash
# 1. Copy environment file
cp .env.example .env

# 2. Edit .env and add your Telegram bot token
nano .env  # or use your favorite editor

# 3. Install air (Go auto-reload)
go install github.com/air-verse/air@latest

# 4. Install dependencies
npm install
cd mini-app && npm install && cd ..
go mod download
```

### Development Mode

```bash
# Start both frontend and backend with one command
npm run dev

# OR start them separately:
npm run dev:server    # Backend only (port 3001)
npm run dev:client    # Frontend only (port 5173)
```

### Build for Production

```bash
# Build everything (frontend + backend binary)
npm run build:all

# Run production binary
./bin/server
```

---

## 🔧 Individual Commands

```bash
# Frontend commands
cd mini-app
npm run dev         # Development server
npm run build       # Production build

# Backend commands
go run ./cmd/server              # Run backend
air                              # Run with auto-reload
go build -o bin/server ./cmd/server  # Build binary

# Database
# SQLite database will be created automatically at ./data/dev.db
```

---

## 🌐 Access URLs

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:3001
- **Health Check**: http://localhost:3001/health
- **API Base**: http://localhost:3001/api/v1

---

## 🤖 Telegram Bot Setup

1. **Create bot** with @BotFather:

   ```
   /newbot
   ```

   Save the token and add it to `.env`

2. **Create Mini App**:

   ```
   /newapp
   ```

   Set URL to your deployment (or use ngrok for local testing)

3. **Set webhook** (optional, for bot commands):
   ```bash
   curl -X POST https://api.telegram.org/bot<YOUR_TOKEN>/setWebhook \
     -d "url=https://your-domain.com/api/v1/telegram/webhook"
   ```

---

## 🧪 Testing

```bash
# Test health endpoint
curl http://localhost:3001/health

# Test auth endpoint (with Telegram initData)
curl -X POST http://localhost:3001/api/v1/auth/telegram/verify \
  -H "Content-Type: application/json" \
  -d '{"initData": "your-telegram-init-data"}'
```

---

## 📦 Project Structure

```
telegram-task-manager/
├── start.sh              # Quick start script
├── cmd/server/           # Go main entry
├── internal/             # Go backend code
├── mini-app/             # Vue 3 frontend
├── .env                  # Environment variables
└── README.md             # Full documentation
```

---

## ⚡ Quick Troubleshooting

**Backend won't start?**

- Check if `.env` exists and has valid values
- Make sure port 3001 is not in use
- Run `go mod download` to ensure dependencies are installed

**Frontend build fails?**

- Delete `mini-app/node_modules` and run `npm install` again
- Make sure you're using Node.js 18+

**Can't connect?**

- Check CORS settings in `.env` (`ALLOWED_ORIGINS`)
- Ensure both servers are running
- Frontend proxy is configured in `mini-app/vite.config.js`




