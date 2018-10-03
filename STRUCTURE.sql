DROP TABLE nodes;
CREATE TABLE nodes (
    id          char(8) PRIMARY KEY,
    account_id  char(8),
    name        varchar(255) NOT NULL,
    url         varchar(255) NOT NULL,
    hostname    varchar(255) NOT NULL,
    port        varchar(255) NOT NULL,
    status      varchar(255) NOT NULL,
    public      boolean,
    created_at  date,
    updated_at  date
);

DROP TABLE accounts;
CREATE TABLE accounts (
    id                        char(8) PRIMARY KEY,
    email                     varchar(255) NOT NULL,
    email_confirmation_token  varchar(255) NOT NULL,
    email_confirmed           boolean,
    username                  varchar(255) NOT NULL,
    ip                        varchar(255) NOT NULL,
    auth_token                varchar(255) NOT NULL,
    created_at                date,
    updated_at                date
);
