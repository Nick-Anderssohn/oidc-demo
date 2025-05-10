-- +goose Up
-- +goose StatementBegin
create table demo.state_token (
    token text primary key,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);

create trigger state_token_updated_at
    before update on demo.state_token
    for each row
    execute procedure set_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table demo.state_token;
-- +goose StatementEnd
