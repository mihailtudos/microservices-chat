-- name: CreateMessage :one
INSERT INTO chat_messages (
    chat_id, from_user, message_content
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: CreateChat :one
INSERT INTO chats (
    name, usernames
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetChatById :one
SELECT * FROM chats
WHERE id = $1;