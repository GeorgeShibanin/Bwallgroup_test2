create table client  (
    id serial primary key,
    balance int
);

create table query  (
     id serial primary key,
     client_id int,
     operation_sum int,
     operation_accepted boolean,
     created_at date,
     foreign key (client_id) references client (id)
);