-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE chirps (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  body VARCHAR(255) NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose statementbegin
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose statementend

CREATE TRIGGER trigger_update_timestamp
BEFORE UPDATE ON chirps
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_update_timestamp ON chirps;
DROP TABLE chirps;
