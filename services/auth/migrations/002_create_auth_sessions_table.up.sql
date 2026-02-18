CREATE TABLE auth_sessions(
  id BIGSERIAL PRIMARY KEY,
  parent_id BIGINT,
  auth_id BIGINT NOT NULL,

  refresh_token VARCHAR UNIQUE NOT NULL,
  user_agent TEXT NOT NULL,
  ip_address TEXT NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ,

  FOREIGN KEY(auth_id) REFERENCES auth(id) ON DELETE CASCADE,
  FOREIGN KEY(parent_id) REFERENCES auth_sessions(id) ON DELETE CASCADE
);

-- Index records by refresh token if active (not revoked)
CREATE INDEX idx_auth_sessions_refresh_token ON auth_sessions(refresh_token) WHERE revoked_at IS NULL;
