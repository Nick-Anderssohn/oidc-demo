-- name: GetUser :one
select *
from demo."user"
where id = $1;

-- name: GetStateToken :one
select *
from demo.state_token
where token = $1;

-- name: InsertStateToken :exec
insert into demo.state_token (token)
values ($1);

-- name: DeleteStateToken :exec
delete from demo.state_token
where token = $1;

-- name: InsertNonce :exec
insert into demo.nonce (nonce)
values ($1);

-- name: UpsertUserByEmail :one
insert into demo."user" (email)
values ($1)
on conflict (email) do update set email = excluded.email
returning *;

-- name: UpsertIdentity :one
insert into demo.identity (user_id, identity_provider_id, external_id, most_recent_id_token)
values ($1, $2, $3, $4)
on conflict (user_id, identity_provider_id, external_id)
do update set most_recent_id_token = excluded.most_recent_id_token
returning *;

-- name: GetSession :one
select *
from demo.session
where id = $1;

-- name: InsertSession :exec
insert into demo.session (id, user_id)
values ($1, $2);

-- name: DeleteSession :exec
delete from demo.session
where id = $1;

-- name: GetUserData :many
select "user".email as user_email,
        identity.id as identity_id,
        identity.identity_provider_id as identity_provider_id,
        identity.external_id as external_id,
        identity.most_recent_id_token as most_recent_id_token
from demo."user"
left join demo.identity identity on identity.user_id = "user".id
where "user".id = $1;

-- name: GetUserByIdentityExternalID :one
select u.*
from demo."user" u
where exists (
        select 1
        from demo.identity i
        where i.user_id = u.id
          and i.identity_provider_id = $1
          and i.external_id = $2
);

-- name: DeleteUser :exec
delete from demo."user" where id = $1;