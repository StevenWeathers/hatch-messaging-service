-- +goose Up
-- +goose StatementBegin
ALTER TABLE hatch.conversation_message ADD COLUMN scheduled_time TIMESTAMPTZ;
ALTER TABLE hatch.conversation_message ADD COLUMN message_sent BOOLEAN;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE hatch.conversation_message DROP COLUMN scheduled_time;
ALTER TABLE hatch.conversation_message DROP COLUMN message_sent;
-- +goose StatementEnd
