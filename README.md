# Revenue & Conversion Intelligence Dashboard (MVP)

Production-grade MVP scaffold for an executive dashboard covering revenue and conversion intelligence.

## What’s Included
- dbt project for core warehouse models (facts + dimensions)
- Go API service skeleton (Fiber) with metrics endpoints
- Next.js TypeScript frontend skeleton with KPI cards
- Docker Compose for local Redis + Airflow placeholder

## Quick Start (Local)

### 1) Start infrastructure
```bash
cd /home/sonthep/dev
docker-compose up -d
```

### 2) API (Go)
```bash
cd /home/sonthep/dev/api
go mod tidy
go run ./
```

The API uses a local SQLite dev warehouse by default. You can override it with:
- `WAREHOUSE_DSN` in [api/.env.example](api/.env.example)

Optional auth:
- Set `API_KEY` in [api/.env.example](api/.env.example)
- Send `X-API-Key: <key>` or `Authorization: Bearer <key>` to access `/api/*`

Sample metric endpoints:
- `/api/metrics/revenue`
- `/api/metrics/conversion-rate`
- `/api/metrics/arpu`
- `/api/metrics/mrr`
- `/api/metrics/nrr`
- `/api/metrics/churn-rate`
- `/api/metrics/ltv`
- `/api/metrics/cac`

### 3) Frontend (Next.js)
```bash
cd /home/sonthep/dev/frontend
npm install
npm run dev
```

### 4) dbt (Data models)
```bash
cd /home/sonthep/dev/data
# Edit profiles.yml with your warehouse credentials
# Then run:
dbt deps
dbt run
```

## Folder Structure
- data/ (dbt models, tests, seeds)
- api/ (Go service)
- frontend/ (Next.js UI)
- orchestration/ (Airflow DAG placeholders)

## Notes
- Update .env.example files and create .env as needed.
- Airflow service is a placeholder; configure as needed for your environment.
- Local SQLite database file is created at api/dev.db on first run.
- Data quality DAG reads SQLite from /opt/airflow/api/dev.db (mounted from ./api).
