-- name: GetUser :one
select *
from demo."user"
where id = $1;