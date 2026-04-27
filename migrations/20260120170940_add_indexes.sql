-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE INDEX IF NOT EXISTS idx_couriers_status_id
ON couriers (status, id);

CREATE INDEX IF NOT EXISTS idx_delivery_courier_id
ON delivery (courier_id);
CREATE INDEX IF NOT EXISTS idx_delivery_courier_deadline
ON delivery (courier_id, deadline);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP INDEX IF EXISTS idx_delivery_courier_deadline;
DROP INDEX IF EXISTS idx_delivery_courier_id;
DROP INDEX IF EXISTS idx_couriers_status_id;
-- +goose StatementEnd
