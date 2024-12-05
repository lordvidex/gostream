-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stream_users (
  id SERIAL PRIMARY KEY NOT NULL,
  name VARCHAR(255) NOT NULL,
  age INT NOT NULL,
  nationality TEXT NOT NULL
)
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE stream_users 
ADD CONSTRAINT IF NOT EXISTS nationality_iso2_chk CHECK(char_length(nationality) = 2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stream_users;
-- +goose StatementEnd
