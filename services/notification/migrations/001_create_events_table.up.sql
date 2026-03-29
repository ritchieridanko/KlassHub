CREATE TABLE events(
  id UUID PRIMARY KEY,
  
  topic TEXT NOT NULL,
  payload JSONB NOT NULL,
  
  first_processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  last_processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  completed_at TIMESTAMPTZ
);
