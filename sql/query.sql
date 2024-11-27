-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserByName :one
SELECT * FROM users
WHERE name = ? LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
	name, password_hash, balance, invites
) VALUES (
	?, ?, ?, ?
)
RETURNING id;

-- name: UpdatePassword :exec
UPDATE users
SET password_hash = ?
WHERE id = ?;

-- name: UpdateBalance :exec
UPDATE users
SET balance = ?
WHERE id = ?;

-- name: GetInvite :one
SELECT * FROM invites
WHERE code = ? LIMIT 1;

-- name: GetUserInvites :many
SELECT * FROM invites
WHERE user_id = ? AND used = FALSE;

-- name: CreateInvite :exec
INSERT INTO invites (
	code, used, user_id
) VALUES (
	?, ?, ?
);

-- name: UpdateUserInvites :exec
UPDATE users
SET invites = invites - 1
WHERE id = ?;

-- name: UseInvite :exec
UPDATE invites
SET used = TRUE
WHERE id = ?;

-- name: ListUserServices :many
SELECT
	services.id, services.name,
	services.expires_at, services.type,
	service_locations.name
FROM services
JOIN service_locations ON service_locations.id = services.location_id
WHERE services.user_id = ?;

-- name: GetService :one
SELECT
	services.id, services.name,
	services.expires_at, services.created_at,
	services.prolong, services.prolong_price,
	services.type,
	service_locations.name,
	service_locations.address
FROM services
JOIN service_locations ON service_locations.id = services.location_id
WHERE services.id = ? AND services.user_id = ?;

-- name: CreateService :one
INSERT INTO services (
	name, type, created_at, expires_at, prolong, prolong_price, user_id, location_id
) VALUES (
	?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING id;

-- name: DeleteService :exec
DELETE FROM services
WHERE id = ?;

-- name: GetExpiredServices :many
SELECT services.*, service_locations.address FROM services
JOIN service_locations ON service_locations.id = services.location_id
WHERE expires_at < ?;

-- name: ProlongService :exec
UPDATE services
SET expires_at = ? + (expires_at - created_at)
WHERE id = ?;

-- name: CreateTransaction :one
INSERT INTO transactions (
	payment_id, amount, status, timestamp, url, user_id
) VALUES (
	?, ?, ?, ?, ?, ?
)
RETURNING id;

-- name: ListTransactions :many
SELECT * FROM transactions
WHERE user_id = ?;

-- name: CancelExpiredTransactions :many
UPDATE transactions
SET status = 'canceled'
WHERE timestamp < ?
AND status = 'in_process'
RETURNING payment_id;

-- name: UpdateTransaction :one
UPDATE transactions
SET status = ?
WHERE payment_id = ?
RETURNING user_id, amount;

-- name: ListLocations :many
SELECT name, services FROM service_locations;

-- name: GetLocation :one
SELECT * FROM service_locations
WHERE services LIKE '%' || ? || '%'
AND name = ? LIMIT 1;

-- name: GetPrice :one
SELECT amount FROM service_prices
WHERE type = ?;
