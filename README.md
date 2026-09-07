# Habit / Task Tracker

Project #1 in the [C# / Go / React skill-building sequence](https://github.com/ChristopherGoodale/csharp-go-react-project-ranking.md):
a CRUD habit/task tracker with a Go backend and a React frontend, backed by
Postgres. No auth, no real-time — the goal is to feel the full
request → API → DB → UI loop end to end with nothing hidden behind a
framework.

- **Backend**: Go, stdlib `net/http` (Go 1.22+ pattern routing), plain SQL
  against Postgres via `database/sql` — no router library, no ORM.
- **Frontend**: React + Vite, plain `fetch` calls, no state management
  library.

See [`habit_task_tracker_guide/`](habit_task_tracker_guide/) for the reasoning
behind these choices.

## Prerequisites

- Go 1.22+
- Node 18+
- Docker Desktop (for local Postgres)

## Run it

1. **Start Postgres:**

   ```bash
   cp .env.example .env
   docker compose up -d
   ```

   This project's Postgres runs on host port **5433** (not the default 5432)
   so it doesn't collide with a Postgres instance you might already have
   running locally. Change `POSTGRES_PORT` and `DATABASE_URL` in `.env` if
   5433 is also taken.

2. **Start the Go server** (reads `.env` values from your shell, or falls
   back to the same defaults as `.env.example`):

   ```bash
   cd server
   go run .
   ```

   The server listens on `http://localhost:8080` and creates the `tasks`
   table on startup if it doesn't already exist.

3. **Start the React client:**

   ```bash
   cd client
   npm install
   npm run dev
   ```

   Open `http://localhost:5173`.

## API

| Method | Path             | Description       |
|--------|------------------|--------------------|
| GET    | `/api/tasks`     | List all tasks     |
| POST   | `/api/tasks`     | Create a task      |
| PUT    | `/api/tasks/{id}`| Update a task      |
| DELETE | `/api/tasks/{id}`| Delete a task      |

## Project layout

```
server/   Go backend (internal/db, internal/task, internal/middleware)
client/   React frontend (Vite)
```
