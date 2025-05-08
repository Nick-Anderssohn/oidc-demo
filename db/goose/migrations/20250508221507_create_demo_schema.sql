-- +goose Up
-- +goose StatementBegin
create schema demo;
GRANT ALL PRIVILEGES ON SCHEMA demo TO demo;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop schema demo;
-- +goose StatementEnd
