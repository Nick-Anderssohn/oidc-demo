-- +goose Up
-- +goose StatementBegin
create table demo.identity_provider (
    id text primary key,

    created_at timestamp default now(),
    updated_at timestamp default now()
);

create trigger identity_provider_set_updated_at
before update on demo.identity_provider
for each row
execute function set_updated_at();

insert into demo.identity_provider (id) values ('google');

create table demo.identity (
    id uuid default uuid_generate_v4() primary key,
    identity_provider_id text not null references demo.identity_provider(id) on delete cascade,
    user_id uuid not null references demo."user"(id) on delete cascade,
    external_id text not null,
    most_recent_id_token jsonb,
    created_at timestamp default now(),
    updated_at timestamp default now(),

    unique (identity_provider_id, external_id, user_id)
);

create trigger identity_set_updated_at
before update on demo.identity
for each row
execute function set_updated_at();

create index idx_identity_user_id on demo.identity(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table demo.identity;
drop table demo.identity_provider;
-- +goose StatementEnd
