-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pets (
  id SERIAL PRIMARY KEY NOT NULL,
  kind VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  age INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pets;
-- +goose StatementEnd
