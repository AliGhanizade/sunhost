# SunHost — Technical Documentation

This document describes the architecture, database schema, API, and design
decisions of SunHost. It is intended both as project documentation and as a
record of the learning process behind the project.

## 1. Overview

SunHost is a small hosting-company website backed by a Go server. It consists
of two main surfaces:

1. **Public website** (`/host/*`) — static Persian/RTL pages rendered from
   server-side HTML templates: home, hosting plans, blog, about, contact, and
   login/registration.
2. **Admin panel** (`/panel/*`) — a dashboard with user-management views and a
   live system-monitoring page fed by Server-Sent Events (SSE).

The backend is a single Go binary built with the Gin framework and SQLite as
the embedded database.

## 2. Architecture

```
Browser
  │  HTML (Go html/template)     │  JSON REST          │  SSE stream
  ▼                              ▼                     ▼
┌─────────────────────────────────────────────────────────────┐
│ router/  (Gin engine, routes, template loading)             │
├─────────────────────────────────────────────────────────────┤
│ controllers/                                                │
│   ├── user_controller         register, login               │
│   ├── log_controller          activity logging              │
│   └── system_stat_controller  metrics stream, summary       │
├─────────────────────────────────────────────────────────────┤
│ model/  (User, UserLog, SystemStat + persistence methods)   │
├─────────────────────────────────────────────────────────────┤
│ config/  (SQLite connection, migrations)                    │
└─────────────────────────────────────────────────────────────┘
                     │
                     ▼
               sunhost.db (SQLite file)
```

The project follows an MVC-like structure:

- `model/` — structs plus the SQL statements that read and write them. Each
  model owns its own queries (a small "active record" style, chosen for
  simplicity in a course project).
- `controllers/` — Gin handlers, one package per feature area.
- `router/` — route registration and template parsing.
- `config/` — database initialization and schema migrations.

## 3. Database Schema

Three tables, created with `CREATE TABLE IF NOT EXISTS` migrations on startup
in `config/config.go`:

### `users`

| column    | type    | constraints         |
| --------- | ------- | ------------------- |
| id        | INTEGER | PK, autoincrement   |
| full_name | TEXT    | NOT NULL            |
| username  | TEXT    | NOT NULL, UNIQUE    |
| email     | TEXT    | NOT NULL, UNIQUE    |
| password  | TEXT    | NOT NULL (bcrypt)   |

Since v1.1.0 passwords are stored as bcrypt hashes.

### `logs`

| column      | type      | constraints                        |
| ----------- | --------- | ---------------------------------- |
| id          | INTEGER   | PK, autoincrement                  |
| username    | TEXT      | NOT NULL, FK → users(username)     |
| ip_address  | TEXT      | NOT NULL                           |
| system_info | TEXT      | NOT NULL (OS + browser string)     |
| action      | TEXT      | NOT NULL                           |
| time        | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP          |

Each user keeps at most 20 log rows: on insert, the controller counts the
user's rows and, once past 20, deletes the oldest non-registration entry. This
keeps the log table bounded without a background job.

### `system_stats`

| column       | type      | constraints       |
| ------------ | --------- | ----------------- |
| id           | INTEGER   | PK, autoincrement |
| alloc_ram    | REAL      | NOT NULL          |
| goroutines   | INTEGER   | NOT NULL          |
| live_objects | INTEGER   | NOT NULL          |
| record_time  | TIMESTAMP | DEFAULT now       |

While the live stream is running, one row is written roughly every 30 seconds
(every 10th tick of the 3-second ticker). A cleanup query keeps only the most
recent 1000 rows so the database cannot grow indefinitely.

Note: SQLite's foreign-key enforcement only works when `PRAGMA foreign_keys =
ON` is executed on the connection, which the config layer does at startup.

## 4. API Reference

### POST /api/users/register

Creates a user.

Request:

```json
{
  "full_name": "Ali Ghanizadeh",
  "username": "ali",
  "email": "ali@example.com",
  "password": "secret"
}
```

Responses:

- `201` — `{ "message": "User created successfully" }`
- `400` — malformed JSON or missing fields
- `409` — `{ "error": "Username already exists" }` or
  `{ "error": "Email already exists" }`
- `500` — database failure

### POST /api/users/login

Verifies credentials. Returns the user object **without** the password field.

Request:

```json
{ "username": "ali", "password": "secret" }
```

Responses:

- `200` — `{ "user": { "id", "full_name", "username", "email" } }`
- `400` — malformed JSON
- `401` — `{ "error": "Invalid username or password" }`

Both "unknown user" and "wrong password" return the same `401` message so the
endpoint does not reveal which usernames exist.

### GET /api/users/logs?username=ali

Returns the stored activity logs of a user:

```json
{ "logs": [ { "id": 1, "username": "ali", "ip_address": "...",
              "system_info": "Windows - Chrome", "action": "User logged in",
              "time": "2026-04-21 19:25:03" } ] }
```

### GET /api/system/summary

Aggregated counters over the `system_stats` table:

```json
{ "avg_ram_today": "12.4", "max_goroutines": 84, "total_records": 530 }
```

### GET /api/system/stream

Server-Sent Events endpoint. Every 3 seconds it emits one `data:` event with a
JSON snapshot:

```json
{
  "cpu_percent": 12.5, "cpu_cores": 8,
  "ram_total": 16384, "ram_used": 7210, "ram_percent": 44.0,
  "swap_total": 4096, "swap_used": 512,
  "disk_total": 476000, "disk_used": 250300, "disk_percent": 52.6,
  "net_sent": 123456789, "net_recv": 987654321, "net_conns": 210,
  "os_name": "windows", "uptime": 85123, "procs": 231,
  "go_alloc_mb": 4.2, "num_gc": 152, "live_objects": 9182,
  "num_goroutine": 12
}
```

Memory/runtime values come from `runtime.MemStats`; machine values come from
`gopsutil` (cpu, mem, disk, net, host).

## 5. Activity Logging

Controllers push context values (`action`, `username`) into the Gin context
and then call `LogController.Create`, which reads them back and inserts a row
with the client IP and a small parsed user-agent string (OS + browser). Log
events include: user registration, successful login, failed login, opening the
live stream, and requesting the summary endpoint.

## 6. Configuration

| Setting | Default        | How to change                  |
| ------- | -------------- | ------------------------------ |
| Port    | `8080`         | `PORT` environment variable    |
| DB path | `./sunhost.db` | Recreate on first run; file is gitignored |

Since v1.1.0 the SQLite connection is limited to a single open connection
(`SetMaxOpenConns(1)`), which avoids `database is locked` errors under
concurrent writes — the standard practice for file-backed SQLite.

## 7. Template Loading

The router parses all `public/*.html` and `public/panel/*.html` files into one
`html/template` set at startup and serves pages with `c.HTML(...)`. Templates
share partials such as `navbar.html` and `footer.html`.

## 8. Frontend Notes

- Bootstrap 5 RTL build is served locally from `public/assets/`.
- `login.js` performs client-side validation (username length, password
  confirmation) before calling the REST API and stores basic session info in
  `localStorage`.
- `panel.js` drives the admin panel: it opens the SSE stream with
  `EventSource`, updates the monitoring cards, and fetches activity logs.
- `products.js` implements client-side plan filtering on the pricing page.

## 9. Security Notes

What version 1.1.0 fixed (2026-09):

- Passwords are hashed with bcrypt instead of being stored in plaintext.
- The login/register responses never include the password field
  (`json:"-"`).
- Login returns a uniform `401` for unknown users and wrong passwords.
- Wrong-password responses were previously empty; they now return JSON.

Known remaining limitations (accepted for a course project):

- No sessions/JWT: the panel trusts `localStorage` set by the login page, so
  "authentication" is client-side only.
- No CSRF protection on the API.
- The activity-log endpoint accepts any username, so any user's logs can be
  read by anyone.
- No rate limiting.

## 10. Development History

The project was developed during the spring 2026 semester for the Web Design
course at Sajjad University:

- **Feb–Mar 2026** — static website (layout, pages, pricing, login form).
- **Mar 2026** — Go backend: Gin server, SQLite layer, models.
- **Mar–Apr 2026** — user API, activity logging, admin panel pages.
- **Apr–May 2026** — live system monitoring with SSE and gopsutil.
- **Jun 2026** — final course submission (`v1.0.0`).
- **Sep 2026** — public release, security refactor (`v1.1.0`).
