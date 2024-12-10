
CREATE TABLE soda_servers
(
    id          int NOT NULL AUTO_INCREMENT,
    name        varchar(255) not null,
    ip_address  varchar(100) not null,
    port        varchar(100) not null,
    username    varchar(100) not null,
    password    varchar(100) not null,
    type        varchar(100) not null,
    PRIMARY KEY (id)
);


CREATE TABLE soda_databases
(
    id           int NOT NULL AUTO_INCREMENT,
    name         varchar(255) not null,
    server_name  varchar(255) not null,
    PRIMARY KEY (id)
);