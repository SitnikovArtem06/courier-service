-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE couriers
    ADD COLUMN transport_type TEXT NOT NULL DEFAULT 'on_foot';


CREATE TABLE delivery (
                          id BIGSERIAL PRIMARY KEY,
                          courier_id BIGINT NOT NULL REFERENCES couriers(id) ON DELETE RESTRICT,
                          order_id VARCHAR(255) NOT NULL UNIQUE,
                          assigned_at TIMESTAMP NOT NULL DEFAULT NOW(),
                          deadline TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
