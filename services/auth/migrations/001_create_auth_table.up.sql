CREATE TABLE auth(
  id BIGSERIAL PRIMARY KEY,

  email VARCHAR,
  username VARCHAR,
  CHECK(email IS NOT NULL OR username IS NOT NULL),

  password VARCHAR NOT NULL, -- hashed
  is_school BOOLEAN NOT NULL,

  last_login_at TIMESTAMPTZ,
  email_verified_at TIMESTAMPTZ,
  email_changed_at TIMESTAMPTZ,
  username_changed_at TIMESTAMPTZ,
  password_changed_at TIMESTAMPTZ,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

-- Enforce uniqueness of email for active (not deleted) records
CREATE UNIQUE INDEX idx_auth_unique_email ON auth(email) WHERE email IS NOT NULL AND deleted_at IS NULL;

-- Enforce uniqueness of username for active (not deleted) records
CREATE UNIQUE INDEX idx_auth_unique_username ON auth(username) WHERE username IS NOT NULL AND deleted_at IS NULL;
