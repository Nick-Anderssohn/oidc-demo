-- +goose Up
-- +goose StatementBegin
create table demo."user" (
    id uuid default uuid_generate_v4() primary key,
    email text not null unique,
    
    created_at timestamp default now(),
    updated_at timestamp default now()
);

create trigger user_set_updated_at
before update on demo."user"
for each row
execute function set_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table demo."user";
-- +goose StatementEnd