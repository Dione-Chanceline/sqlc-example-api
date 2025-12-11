--migrate Up



CREATE TABLE "message" (
  "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
  "thread" VARCHAR(36) NOT NULL,
  "sender" VARCHAR(100) NOT NULL,
  "content" TEXT NOT NULL,
  "created_at" TIMESTAMP DEFAULT now()
);


CREATE TABLE attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
    message_id VARCHAR(36) NOT NULL REFERENCES message(id) ON DELETE CASCADE,
    file_url TEXT NOT NULL,
    file_type TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

