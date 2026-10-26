POST save-operation
body
id_client int
id_external int // внешний идентификатор операции
summ int // размер операции в копейках
time timestamptz // время операции

Postgres
table name operation_history
id BIGSERIAL PRIMARY KEY
client_id int
id_external int UNIQUE
summ int
time timestamptz


GET operation-history/{client_id}?id_op=id

select *
from operation_history
where client_id = $1
offset $2
limit $3

select *
from operation_history
where client_id = $1 and id > $2
order by time
limit $3