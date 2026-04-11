CREATE TABLE user_preferences(
  auth_id BIGINT PRIMARY KEY,

  theme VARCHAR NOT NULL DEFAULT 'system',
  language VARCHAR NOT NULL DEFAULT 'en',
  timezone VARCHAR NOT NULL DEFAULT 'asia/jakarta',

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);
