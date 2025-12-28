# SpotiStats

Music analytics backend with Go, Kafka, Redis, and Postgres.

## Features

- REST API with Swagger docs
- 50k Spotify tracks dataset
- Redis caching
- Kafka message queue
- Background job processing
- Real-time dashboard

## Quick Start
```bash
# Clone
git clone https://github.com/yourusername/spotistats.git
cd spotistats

# Setup env
cp .env.example .env
# Edit .env and set DB_PASSWORD

# Run (migrate + seed + serve)
docker compose up --build -d

# Open
open http://localhost:8080
```

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Dashboard |
| GET | `/health` | Health check |
| GET | `/tracks` | List tracks |
| GET | `/tracks/:id` | Get track |
| POST | `/tracks` | Add track (via Kafka) |
| POST | `/jobs` | Create background job |
| GET | `/jobs/:id` | Get job status |
| GET | `/swagger/*` | API docs |

## Tech Stack

- **API:** Go + Echo
- **Database:** PostgreSQL
- **Cache:** Redis
- **Queue:** Apache Kafka
- **Docs:** Swagger
