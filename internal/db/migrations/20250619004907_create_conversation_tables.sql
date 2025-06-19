-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS hatch.contact (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone_number VARCHAR(20),
    email VARCHAR(255)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_contact_phone_number ON hatch.contact(phone_number); -- Ensure phone number is unique
CREATE UNIQUE INDEX IF NOT EXISTS idx_contact_email ON hatch.contact(LOWER(email)); -- Ensure email is unique and case-insensitive

CREATE TABLE IF NOT EXISTS hatch.conversation (
    id SERIAL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS hatch.conversation_contact (
    conversation_id INT NOT NULL REFERENCES hatch.conversation(id) ON DELETE CASCADE,
    contact_id INT NOT NULL REFERENCES hatch.contact(id) ON DELETE CASCADE,
    PRIMARY KEY (conversation_id, contact_id)
);

CREATE TABLE IF NOT EXISTS hatch.conversation_message (
    id SERIAL PRIMARY KEY,
    conversation_id INT NOT NULL REFERENCES hatch.conversation(id) ON DELETE CASCADE,
    sender_id INT NOT NULL REFERENCES hatch.contact(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    attachments JSONB DEFAULT '[]',
    type VARCHAR(20) NOT NULL CHECK (type IN ('email', 'sms', 'mms')),
    timestamp TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_conversation_message_type ON hatch.conversation_message(type);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_contact_phone_number;
DROP INDEX IF EXISTS idx_contact_email;
DROP INDEX IF EXISTS idx_conversation_message_type;
DROP TABLE IF EXISTS hatch.conversation_contact;
DROP TABLE IF EXISTS hatch.conversation_message;
DROP TABLE IF EXISTS hatch.conversation;
DROP TABLE IF EXISTS hatch.contact;
-- +goose StatementEnd
