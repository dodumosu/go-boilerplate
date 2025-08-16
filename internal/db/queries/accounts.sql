-- name: CreateUser :one
-- Create a new user
INSERT INTO users (
    id,
    email,
    username,
    password_hash,
    is_verified
) VALUES (
    sqlc.arg(id),
    sqlc.arg(email),
    sqlc.arg(username),
    sqlc.arg(password_hash),
    NOT sqlc.arg(requires_verification)::boolean
)
RETURNING *;

-- name: CreateProfile :one
-- Create a new profile for a user
INSERT INTO profiles (
    id,
    user_id,
    first_name,
    last_name,
    other_names,
    bio,
    phone
) VALUES (
    sqlc.arg(id),
    sqlc.arg(user_id),
    sqlc.arg(first_name),
    sqlc.arg(last_name),
    sqlc.arg(other_names),
    sqlc.arg(bio),
    sqlc.arg(phone)
)
RETURNING *;

-- name: CreateOAuthAccount :one
-- Create a new oauth account for a user
INSERT INTO oauth_accounts (
    id,
    user_id,
    provider_id,
    provider_user_id,
    provider_username,
    provider_email,
    access_token,
    refresh_token,
    token_expires_at,
    raw_user_data
) VALUES (
    sqlc.arg(id),
    sqlc.arg(user_id),
    sqlc.arg(provider_id),
    sqlc.arg(provider_user_id),
    sqlc.arg(provider_username),
    sqlc.arg(provider_email),
    sqlc.arg(access_token),
    sqlc.arg(refresh_token),
    sqlc.arg(token_expires_at),
    sqlc.arg(raw_user_data)
)
RETURNING *;
