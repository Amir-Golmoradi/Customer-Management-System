# CustomerSystem Runbook

## Start
- `docker compose up --build`

## Endpoints
- `GET /health/live`
- `GET /health/ready`
- `GET /metrics`
- `POST /v1/customers`
- `GET /v1/customers?limit=20&offset=0`
- `GET /v1/customers/{id}`

## Graceful restart
1. Send `SIGTERM` to API process/container.
2. Wait for `HTTP_SHUTDOWN_TIMEOUT` drain window.
3. Confirm `/health/live` is down and replacement instance is ready.

## Reliability checks
- `/health/ready` must return `200` before adding instance to traffic.
- Track `http_request_errors_total` and `db_query_errors_total`.
- Watch pool saturation via `db_pool_acquired_conns` vs `db_pool_max_conns`.

## Backup/restore checklist
- Enable daily Postgres logical dump.
- Keep at least 7 backup snapshots.
- Test restore monthly in isolated environment.
