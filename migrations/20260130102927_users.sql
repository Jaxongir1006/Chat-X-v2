-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        CREATE TYPE user_role AS ENUM ('user', 'admin', 'superuser');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE,
    phone VARCHAR(20) UNIQUE,
    email VARCHAR(200) UNIQUE NOT NULL,
    password_hash VARCHAR(500) NOT NULL,
    role user_role DEFAULT 'user',
    verified BOOLEAN DEFAULT FALSE,
    last_seen_at TIMESTAMPTZ,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_profile (
  user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  fullname VARCHAR(250) NOT NULL DEFAULT '',
  bio      VARCHAR(500) NOT NULL DEFAULT '',
  address  VARCHAR(500) NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS user_profile_images (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    image_key VARCHAR(255) NOT NULL, 

    is_primary BOOLEAN DEFAULT FALSE, -- So the frontend knows which one is the "main" profile pic
    display_order INT DEFAULT 0,      -- So the user can drag-and-drop rearrange their photos
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);


CREATE INDEX idx_user_email ON users (email);
CREATE INDEX idx_user_profile_user_id ON user_profile (user_id);
CREATE INDEX idx_user_profile_images_user_id ON user_profile_images (user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_profile;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS user_profile_images;

DROP TYPE IF EXISTS user_role;

DROP INDEX IF EXISTS idx_user_email;
DROP INDEX IF EXISTS idx_user_profile_user_id;
-- +goose StatementEnd
