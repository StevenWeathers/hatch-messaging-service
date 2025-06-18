-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS hatch;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS hatch CASCADE;
-- +goose StatementEnd
