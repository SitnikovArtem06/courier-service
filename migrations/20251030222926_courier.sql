-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE couriers (
                          id          BIGSERIAL PRIMARY KEY,
                          name        TEXT NOT NULL,
                          phone       TEXT NOT NULL UNIQUE,
                          status      TEXT NOT NULL, -- например: 'available', 'busy', 'paused'
                          created_at  TIMESTAMP DEFAULT now(),
                          updated_at  TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
