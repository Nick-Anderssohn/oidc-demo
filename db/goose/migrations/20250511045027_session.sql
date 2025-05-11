-- +goose Up
-- +goose StatementBegin
create table demo.session (
    id text primary key,
    user_id uuid not null references demo."user"(id) on delete cascade,
    created_at timestamp default now(),
    updated_at timestamp default now()
);
create trigger session_set_updated_at
before update on demo.session
for each row
execute function set_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table demo.session;
-- +goose StatementEnd
