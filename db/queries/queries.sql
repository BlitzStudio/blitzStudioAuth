-- name: FindAllUsers :many
SELECT
    *
FROM
    `users`;

-- name: FindUserByEmail :one
SELECT
    *
FROM
    `users`
WHERE
    `email` = ?;

-- name: CountUserByEmail :one
SELECT
    count(*)
FROM
    `users`
WHERE
    `email` = ?;

-- name: CreateUser :execresult
INSERT INTO
    `users` (`email`, `name`, `password`)
VALUES
    (
        sqlc.arg ("email"),
        sqlc.arg ("name"),
        sqlc.arg ("password")
    );

-- name: GetLastInsertedId :one
SELECT
    LAST_INSERTED_ID ();

-- name: FindUserById :one
SELECT
    *
FROM
    `users`
WHERE
    `id` = ?;

-- name: SaveRefreshToken :exec
INSERT INTO
    `jwts` (
        `id`,
        `userId`,
        `tokenFamily`,
        -- `tokenHash`,
        `expiresAt`
    )
VALUES
    (
        sqlc.arg ("tokenId"),
        sqlc.arg ("userId"),
        sqlc.arg ("tokenFamily"),
        -- sqlc.arg ("tokenHash"),
        sqlc.arg ("expiresAt")
    );

-- name: FindRefreshTokenByUserId :one
SELECT
    *
FROM
    `jwts`
WHERE
    `userId` = sqlc.arg ("userId");

-- name: FindValidRefreshTokenByFamilyAndUserId :one
SELECT
    *
FROM
    `jwts`
WHERE
    `tokenFamily` = sqlc.arg ("tokenFamily")
    AND `userId` = sqlc.arg ("userId")
    AND `isRevoked` = FALSE;

-- name: FindTokenById :one
SELECT
    *
FROM
    `jwts`
WHERE
    `id` = sqlc.arg ("id");

-- name: RevokeRefreshTokenById :exec
UPDATE `jwts`
SET
    `isRevoked` = TRUE
WHERE
    `id` = sqlc.arg ("tokenId");

-- name: RevokeTokenFamily :exec
UPDATE `jwts`
SET
    `isRevoked` = TRUE
WHERE
    `tokenFamily` = sqlc.arg ("tokenFamily");