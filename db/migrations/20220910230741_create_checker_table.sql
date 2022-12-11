-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS checker;

CREATE TABLE IF NOT EXISTS checker.users (
    id bigint UNIQUE
);

CREATE TABLE IF NOT EXISTS checker.income_info (
    id            bigint UNIQUE,
    user_name     text,
    bot_link      text,
    bot_name      text,
    income_source text,
    type_bot      text
);

CREATE TABLE IF NOT EXISTS checker.url (
    url_text text,
    url  text
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA checker CASCADE;
-- +goose StatementEnd
