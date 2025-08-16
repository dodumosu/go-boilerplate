-- name: LookupAccountForAuth :one
-- Retrieves an account for the purpose of authentication
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
    (u.username = sqlc.arg(login_identifier) OR u.email = sqlc.arg(login_identifier));
-- AND
--     u.is_verified = TRUE
-- AND
--     u.is_active = TRUE
-- AND
--     (u.is_banned IS NULL OR u.is_banned = FALSE)
-- AND
--     (u.suspended_until IS NULL OR u.suspended_until <= CURRENT_TIMESTAMP)
-- AND
--     (u.deactivate_at IS NULL OR u.deactivate_at > CURRENT_TIMESTAMP);
