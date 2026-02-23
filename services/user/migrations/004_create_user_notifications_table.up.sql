CREATE TABLE user_notifications(
  id UUID PRIMARY KEY,
  auth_id BIGINT NOT NULL,

  type VARCHAR NOT NULL, -- e.g. new thread, new class, new assignment, etc.
  reference_id UUID,
  title VARCHAR NOT NULL,
  content TEXT,
  read_at TIMESTAMPTZ,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index records by auth id
CREATE INDEX idx_user_notifications_auth ON user_notifications(auth_id);
