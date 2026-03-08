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
    `users` (`email`, `name`, `password`, `refresh_token`)
VALUES
    (?, ?, ?, ?);

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