create table client  (
                        id int primary key auto_increment,
                        balance int
);

create table query  (
                         id int primary key auto_increment,
                         id_client int
                         operation_sum int
                         operation_accepted boolean
                         created_at   date
                         foreign key (id_client) references client (id)
);