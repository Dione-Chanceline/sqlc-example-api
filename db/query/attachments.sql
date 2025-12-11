

-- name: InsertAttachment :one
INSERT INTO attachments (message_id, file_url, file_type)
VALUES ($1, $2, $3)
RETURNING id, message_id, file_url, file_type, created_at;

-- name: GetAttachmentsByMessage :many
SELECT id, message_id, file_url, file_type, created_at
FROM attachments
WHERE message_id = $1;


