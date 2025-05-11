-- +goose Up
-- +goose StatementBegin
create table demo.nonce (
    nonce text primary key,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);

create trigger nonce_updated_at
    before update on demo.nonce
    for each row
    execute procedure set_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table demo.nonce;
-- +goose StatementEnd
