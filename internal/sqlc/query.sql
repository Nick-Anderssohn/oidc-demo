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
insert into demo.identity (user_id, identity_provider_id, external_id)
values ($1, $2, $3)
on conflict (user_id, identity_provider_id, external_id)
do update set user_id = excluded.user_id
returning *;