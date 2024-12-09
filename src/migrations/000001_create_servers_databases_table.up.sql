
CREATE TABLE soda_servers
(
    id          int,
    name        varchar(255),
    ip_address  varchar(100)
);


CREATE TABLE soda_databases
(
    id           int,
    name         varchar(255),
    server_name  varchar(255)
);