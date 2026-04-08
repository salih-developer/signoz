# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

SigNoz is an open-source observability platform (APM, distributed tracing, logs, metrics, alerts). It has a Go backend with ClickHouse storage and a React/TypeScript frontend.

## Repository Structure

- `frontend/` — React + Vite + TypeScript frontend
- `cmd/community/` — Community edition Go binary entrypoint
- `cmd/enterprise/` — Enterprise edition Go binary entrypoint
- `pkg/` — Shared Go packages (the bulk of backend logic)
- `ee/` — Enterprise-only Go packages
- `pkg/query-service/` — Legacy query service (excluded from linting, being replaced)
- `tests/integration/` — Python integration tests
- `docs/api/openapi.yml` — OpenAPI 3.0 spec for the backend API
- `.devenv/docker/` — Docker Compose files for local dev services

## Development Commands

### Backend (Go 1.25)

```bash
make devenv-up              # Start ClickHouse + OTel Collector (required first)
make go-run-enterprise      # Run enterprise server (with race detector)
make go-run-community       # Run community server (with race detector)
make go-test                # Run all Go tests with -race
make go-build-enterprise    # Build enterprise binary
make go-build-community     # Build community binary
make gen-mocks              # Generate mocks via mockery (.mockery.yml)
```

### Frontend (Node >= 22, Yarn)

```bash
cd frontend
yarn install
yarn dev                    # Vite dev server on http://localhost:3301
yarn build                  # Production build
yarn lint                   # ESLint
yarn lint:fix               # ESLint auto-fix
yarn prettify               # Prettier auto-format
yarn jest                   # Run tests
yarn jest:watch             # Watch mode
yarn jest:coverage          # With coverage
yarn generate:api           # Generate API types from OpenAPI spec (Orval)
```

### Integration Tests (Python, uv)

```bash
make py-test                # Run all integration tests
make py-test-setup          # Bootstrap test environment
make py-test-teardown       # Teardown test environment
make py-fmt                 # Format with black
make py-lint                # Lint with isort, autoflake, pylint
```

## Architecture

### Backend

The backend uses a modular architecture. The `pkg/signoz` package wires everything together. Key subsystems:

- **`pkg/apiserver`** — REST API (Gorilla mux), OpenAPI-driven, versioned at `/api/v1/` through `/api/v5/`
- **`pkg/telemetrystore`** — Storage abstraction with ClickHouse backend
- **`pkg/telemetrytraces`, `telemetrylogs`, `telemetrymetrics`** — Signal-specific query/storage logic
- **`pkg/alertmanager`, `pkg/ruler`** — Alert management and rule evaluation
- **`pkg/authz`** — Authorization via OpenFGA
- **`pkg/authn`, `pkg/tokenizer`** — Authentication and JWT tokens
- **`pkg/sqlstore`** — SQLite/PostgreSQL for application metadata
- **`pkg/sqlmigration`** — Database migrations
- **`pkg/querier`, `pkg/querybuilder`** — Query execution and construction
- **`pkg/errors`** — Custom errors package (must use instead of stdlib `errors`)
- **`pkg/instrumentation`** — Structured logging via slog

### Frontend

- **State management**: React Query (server state), Zustand (global client state), nuqs (URL state). Do NOT use Redux or Context for new code.
- **UI**: Ant Design 5 + custom SigNoz design system (`@signozhq/ui`, `@signozhq/icons`, `@signozhq/table`)
- **API layer**: Generated React Query hooks from OpenAPI spec via Orval (`frontend/src/api/generated/`)
- **Visualization**: Chart.js, Uplot, D3, Visx
- **Testing**: Jest + Testing Library (80% line coverage threshold)
- **E2E**: Playwright (`frontend/e2e/tests/`)

## Code Conventions

### Go

- **Errors**: Use `github.com/SigNoz/signoz/pkg/errors`, NOT stdlib `errors` or `fmt.Errorf`
- **Logging**: Use `slog` (structured, via `pkg/instrumentation`), NOT `zap` or `fmt.Print*`
- **slog rules**: `attr-only: true`, `context: all`, `static-msg: true`, `key-naming-case: snake`
- **Linter**: golangci-lint with config in `.golangci.yml`. Excluded paths: `pkg/query-service`, `ee/query-service`, `scripts/`, `third_party`

### Frontend

- **Formatting**: Prettier — tabs, single quotes, trailing commas, semicolons, 80 char width, LF line endings
- **Linting**: ESLint with TypeScript strict, React Hooks, SonarJS, simple-import-sort
- **Commits**: Conventional commits enforced via CommitLint + Husky pre-commit hooks
- **Zustand**: Always use selectors, never mutate state directly, one store per module
- **React Query**: Prefer generated hooks from `frontend/src/api/generated/` when available

### Environment

Backend dev server expects ClickHouse at `tcp://127.0.0.1:9000` and uses SQLite (`signoz.db`) by default. Frontend dev server connects to backend at `http://localhost:8080` (configurable via `VITE_FRONTEND_API_ENDPOINT`).
