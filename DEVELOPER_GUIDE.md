# Developer Guide: Telegram Task Manager

Welcome to the **Telegram Task Manager** project! This guide is designed to help new developers quickly understand the architecture, how to navigate the codebase, and how to start making updates.

## 🏗 System Architecture

This application is a full-stack web application designed to be run as a Telegram Mini App.

### Backend (Go)
- **Framework**: Gin (HTTP router)
- **Database**: SQLite (Development) / PostgreSQL (Production) using GORM.
- **Location**: `internal/` directory.
- **Entry Point**: `cmd/server/main.go`
- **Key Features**:
  - `internal/handlers`: Contains all the REST API controllers (Projects, Tasks, Goals, etc.).
  - `internal/services`: Contains the core business logic, including Telegram Bot API interactions (`telegram.go`), background jobs (`scheduler.go`), and Firebase Auth verification (`firebase_auth.go`).
  - `internal/models`: Defines the database schemas.

### Frontend (Vue 3 + Vite)
- **Framework**: Vue 3 (Composition API) with Vite.
- **State Management**: Pinia (`mini-app/src/stores/`).
- **Styling**: Tailwind CSS.
- **Location**: `mini-app/` directory.
- **Key Features**:
  - `src/pages`: Contains the main views (Dashboard, Tasks, Settings, Login).
  - `src/services/api.js`: An Axios instance that automatically attaches the JWT token and handles base URL configuration.
  - `src/utils/firebase.js`: Handles Google and Email authentication via the Firebase SDK.

## 🔑 Authentication Flow

The application supports three types of authentication:
1. **Telegram WebApp**: When opened inside Telegram, the mini-app receives `initData`. The frontend sends this to `POST /api/v1/auth/telegram/verify`, the backend validates the hash, and returns a JWT.
2. **Telegram Widget**: For desktop web users, they can log in using the Telegram Login Widget, which sends data to `POST /api/v1/auth/telegram/widget`.
3. **Firebase (Email/Google)**: Users can log in manually using Email/Password or Google. The frontend generates a Firebase ID Token, sends it to `POST /api/v1/auth/firebase`, and the backend creates/links the user and returns our custom JWT.

## 📝 What Needs to be Done Next?

Based on the current project roadmap and documentation, the following features are still pending implementation:

1. **Reports and Analytics**: We track tasks, habits, and activities, but we need pages and API endpoints to generate charts and insights (e.g., weekly productivity scores).
2. **E2E Tests**: The project needs End-to-End testing (e.g., using Cypress or Playwright) to guarantee frontend flows work without breaking.
3. **Production Deployment Scripts**: We need Docker Compose files, GitHub Actions, or deployment scripts to easily push this to production (e.g., AWS, DigitalOcean).
4. **Mini-App README**: The `mini-app/README.md` is still the default Vite template and needs to be updated with frontend-specific setup instructions.

## 🛠 How to Make an Update

### 1. Adding a New API Endpoint
1. Define the database model in `internal/models/`.
2. Run database migrations in `cmd/server/main.go` by adding your model to `db.AutoMigrate()`.
3. Create a new handler in `internal/handlers/`.
4. Register the route in `cmd/server/main.go` inside the `setupRoutes` function.

### 2. Updating the Frontend UI
1. Create or modify a component in `mini-app/src/components/` or a page in `mini-app/src/pages/`.
2. If it requires global state or API calls, add the logic to a Pinia store in `mini-app/src/stores/`.
3. Ensure you use Tailwind CSS for styling to maintain consistency.

### 3. Local Development
Always use the provided `start.sh` script to run both the Go backend and Vue frontend concurrently:
```bash
./start.sh
```
This enables hot-reloading for both the frontend (Vite) and backend (Air).

## 📚 Reference Documents
- `README.md`: High-level overview and environment variable configurations.
- `QUICK_START.md`: Commands for running the application.
- `notifications_spec.md`: Contains the specifications and UI designs for Telegram Bot notifications.
