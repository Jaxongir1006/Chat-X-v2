-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'conversation_type') THEN
    CREATE TYPE conversation_type AS ENUM ('dm', 'group', 'channel');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'participant_role') THEN
    CREATE TYPE participant_role AS ENUM ('owner', 'admin', 'member', 'restricted', 'banned', 'left');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'message_type') THEN
    CREATE TYPE message_type AS ENUM ('text', 'photo', 'video', 'file', 'voice', 'sticker', 'system');
  END IF;
END
$$;

CREATE TABLE IF NOT EXISTS conversations (
  id            BIGSERIAL PRIMARY KEY,
  type          conversation_type NOT NULL,

  title         TEXT,
  username      TEXT UNIQUE,
  description   TEXT,

  is_public     BOOLEAN NOT NULL DEFAULT FALSE,

  created_by    BIGINT NOT NULL,       -- FK -> users(id)
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),

  -- in order to see the last message fast
  last_message_id BIGINT
);

CREATE TABLE IF NOT EXISTS conversation_participants (
  conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
  user_id         BIGINT NOT NULL, -- FK -> users(id)

  role            participant_role NOT NULL DEFAULT 'member',

  joined_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  left_at         TIMESTAMPTZ,
  muted_until     TIMESTAMPTZ,
  is_pinned       BOOLEAN NOT NULL DEFAULT FALSE,

  last_read_message_id BIGINT,        -- for unread counts
  last_read_at         TIMESTAMPTZ,

  PRIMARY KEY (conversation_id, user_id)
);

CREATE TABLE IF NOT EXISTS dm_pairs (
  user1_id BIGINT NOT NULL,
  user2_id BIGINT NOT NULL,
  conversation_id BIGINT NOT NULL UNIQUE REFERENCES conversations(id) ON DELETE CASCADE,
  PRIMARY KEY (user1_id, user2_id),
  CHECK (user1_id < user2_id)
);

CREATE TABLE IF NOT EXISTS messages (
  id              BIGSERIAL PRIMARY KEY,
  conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,

  sender_id       BIGINT,  -- FK -> users(id), nullable for system messages
  type            message_type NOT NULL DEFAULT 'text',

  text            TEXT,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  edited_at       TIMESTAMPTZ,

  reply_to_id     BIGINT REFERENCES messages(id) ON DELETE SET NULL,
  forward_from_id BIGINT REFERENCES messages(id) ON DELETE SET NULL,

  -- Soft delete -- messages can be restored
  deleted_at      TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS message_reads (
  message_id BIGINT NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
  user_id    BIGINT NOT NULL, -- FK -> users(id)
  read_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (message_id, user_id)
);

CREATE TABLE IF NOT EXISTS message_reactions (
  message_id BIGINT NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
  user_id    BIGINT NOT NULL,
  emoji      TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (message_id, user_id, emoji)
);

CREATE TABLE IF NOT EXISTS message_attachments (
  id          BIGSERIAL PRIMARY KEY,
  message_id  BIGINT NOT NULL REFERENCES messages(id) ON DELETE CASCADE,

  storage_key TEXT NOT NULL,
  url         TEXT,
  mime_type   TEXT,
  size_bytes  BIGINT,
  width       INT,
  height      INT,

  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS pinned_messages (
  conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
  message_id      BIGINT NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
  pinned_by       BIGINT NOT NULL,
  pinned_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (conversation_id, message_id)
);

CREATE TABLE IF NOT EXISTS invite_links (
  id              BIGSERIAL PRIMARY KEY,
  conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,

  code            TEXT NOT NULL UNIQUE,      -- random token
  created_by      BIGINT NOT NULL,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

  expires_at      TIMESTAMPTZ,
  max_uses        INT,
  uses_count      INT NOT NULL DEFAULT 0,

  -- if true: anyone can join
  is_active       BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS join_requests (
  conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
  user_id         BIGINT NOT NULL,
  requested_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  message         TEXT,
  PRIMARY KEY (conversation_id, user_id)
);


CREATE INDEX IF NOT EXISTS idx_invites_conv ON invite_links(conversation_id);
CREATE INDEX IF NOT EXISTS idx_attachments_msg ON message_attachments(message_id);
CREATE INDEX IF NOT EXISTS idx_reactions_msg ON message_reactions(message_id);
CREATE INDEX IF NOT EXISTS idx_messages_conv_created ON messages(conversation_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_messages_reply ON messages(reply_to_id);
CREATE INDEX IF NOT EXISTS idx_participants_user ON conversation_participants(user_id);
CREATE INDEX IF NOT EXISTS idx_participants_conv ON conversation_participants(conversation_id);
CREATE INDEX IF NOT EXISTS idx_conversations_type ON conversations(type);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS conversation_participants;j
DROP TABLE IF EXISTS dm_pairs;
DROP TABLE IF EXISTS message_reads;
DROP TABLE IF EXISTS message_attachments;
DROP TABLE IF EXISTS pinned_messages;
DROP TABLE IF EXISTS invite_links;
DROP TABLE IF EXISTS join_requests;
DROP TABLE IF EXISTS message_reactions;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS conversations;

DROP TYPE IF EXISTS conversation_type;
DROP TYPE IF EXISTS participant_role;
DROP TYPE IF EXISTS message_type;

DROP INDEX IF EXISTS idx_invites_conv;
DROP INDEX IF EXISTS idx_attachments_msg;
DROP INDEX IF EXISTS idx_reactions_msg;
DROP INDEX IF EXISTS idx_messages_conv_created;
DROP INDEX IF EXISTS idx_messages_reply;
DROP INDEX IF EXISTS idx_participants_user;
DROP INDEX IF EXISTS idx_participants_conv;
DROP INDEX IF EXISTS idx_conversations_type;

-- +goose StatementBegin

-- +goose StatementEnd
