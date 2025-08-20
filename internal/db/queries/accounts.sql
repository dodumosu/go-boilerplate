-- name: CheckUserExistsByUsernameOrEmail :one
-- Checks if a user exists with the given username or email.
-- Returns a simple struct indicating existence for each.
SELECT EXISTS(
        SELECT 1
        FROM users
        WHERE users.username = sqlc.arg(username_val)
    ) AS username_exists,
    EXISTS(
        SELECT 1
        FROM users
        WHERE users.email = sqlc.arg(email_val)
    ) AS email_exists;

-- name: CreateUser :one
-- Create a new user
INSERT INTO users (
    id,
    email,
    username,
    password_hash,
    is_active,
    is_verified,
    password_change_on_login,
    is_superuser
) VALUES (
    sqlc.arg(id),
    sqlc.arg(email),
    sqlc.arg(username),
    sqlc.arg(password_hash),
    sqlc.arg(is_active),
    NOT sqlc.arg(requires_verification)::boolean,
    sqlc.arg(password_change_on_login),
    sqlc.arg(is_superuser)
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

-- name: VerifyUserAccount :exec
-- Marks a user's account as verified.
UPDATE users
SET 
    is_verified = TRUE
WHERE id = sqlc.arg(id);


-- name: UpdatePassword :exec
-- Updates a user's password hash
UPDATE users
SET password_hash = sqlc.arg(password_hash)
WHERE id = sqlc.arg(id);

-- name: GetUserByEmail :one
-- Retrieves a user by their email.
SELECT * FROM users WHERE email = sqlc.arg(email) LIMIT 1;


-- name: UpdateUserProfile :one
-- Updates a user's profile information.
UPDATE profiles
SET 
    first_name = COALESCE(sqlc.narg(first_name), first_name),
    last_name = COALESCE(sqlc.narg(last_name), last_name),
    other_names = COALESCE(sqlc.narg(other_names), other_names),
    phone = COALESCE(sqlc.narg(phone), phone),
    bio = COALESCE(sqlc.narg(bio), bio)
WHERE user_id = sqlc.arg(user_id)
RETURNING *;

-- name: GetUserByID :one
-- Retrieves a user by the ID
SELECT * FROM users WHERE id = sqlc.arg(id) LIMIT 1;

-- name: SetPasswordResetFlag :exec
-- Sets the password reset flag on an account
UPDATE users
SET password_reset_requested = true
WHERE email = sqlc.arg(email);


-- name: ResetPassword :exec
-- Reset's an account's password
UPDATE users
SET
    password_hash = sqlc.arg(password_hash),
    password_reset_requested = false
WHERE
    id = sqlc.arg(id);

-- name: GetAccountByUserID :one
-- Retrieves the full account by ID
SELECT
    u.id,
    u.email,
    u.username,
    u.password_hash,
    u.is_superuser,
    u.has_roles,
    u.is_active,
    u.is_verified,
    u.password_reset_requested,
    u.password_change_on_login,
    u.suspended_until,
    u.banned_at,
    u.deactivate_at,
    u.created_at,
    u.updated_at,
    p.id as profile_id,
    p.user_id,
    p.first_name,
    p.last_name,
    p.other_names,
    p.bio,
    p.phone,
    p.created_at as profile_created_at,
    p.updated_at as profile_updated_at
FROM users AS u
INNER JOIN profiles AS p
ON u.id = p.user_id
WHERE
    u.id = sqlc.arg(user_id);