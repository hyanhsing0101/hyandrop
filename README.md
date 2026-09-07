# HyanDrop

HyanDrop is a browser-based temporary transfer tool for moving text, links, images, screenshots, and files across devices through short-lived rooms.

## Stack

- Backend: Go, Gin, WebSocket-ready architecture
- Database: PostgreSQL
- Cache/state: Redis
- Frontend: React, Vite, TypeScript
- Runtime: Docker Compose

## Development

Start the local stack:

```bash
docker compose up --build
```

Default local endpoints:

- Frontend: <http://localhost:5173>
- Backend health: <http://localhost:18080/healthz>
- Backend readiness: <http://localhost:18080/readyz>

The development Compose file starts PostgreSQL, Redis, the Go API, and the Vite frontend. Source code is mounted into the backend and frontend containers, while database and upload data are persisted through Docker volumes or local runtime directories.

## Environment

Runtime defaults are defined in `docker-compose.yml`. Copy `.env.example` to `.env` when you need to override local values.

Do not commit real secrets, production environment files, uploaded files, database data, or logs.

## Project Scope

Initial MVP:

- Create and join temporary rooms
- Join by room code or QR code
- Realtime text/link transfer
- 100 MB file upload and download
- Room history visible after reconnect
- Rooms and files expire 24 hours after the latest successful transfer

Planned later:

- Chunked upload
- Resume interrupted uploads
- Redis-backed rate limiting
- Redis Pub/Sub for multi-instance WebSocket fanout
- Server deployment behind Nginx

## Design

The confirmed first-release feature and privacy design is documented in [docs/functional-design.md](docs/functional-design.md).
