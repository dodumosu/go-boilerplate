-- name: CreateEvent :one
-- Creates an event
INSERT INTO events (id, name, description)
VALUES (sqlc.arg(id), sqlc.arg(name), sqlc.narg(description))
RETURNING *;

-- name: RecordMetric :exec
-- Records a metric
INSERT INTO metrics (
    id,
    event_id,
    dimensions,
    "date",
    "count"
) VALUES (
    sqlc.arg(id),
    sqlc.arg(event_id),
    sqlc.arg(date)::date,
    sqlc.arg(dimensions)::jsonb,
    sqlc.arg(count)
);
