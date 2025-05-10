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