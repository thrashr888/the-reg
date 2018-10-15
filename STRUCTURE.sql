DROP TABLE nodes;
CREATE TABLE nodes (
    id          char(12) PRIMARY KEY,
    account_id  char(12),
    name        varchar(255) NOT NULL,
    url         varchar(255) NOT NULL,
    hostname    varchar(255) NOT NULL,
    port        varchar(255) NOT NULL,
    status      varchar(255) NOT NULL,
    public      boolean,
    created_at  timestamp,
    updated_at  timestamp
);

DROP TABLE accounts;
CREATE TABLE accounts (
    id                        char(12) PRIMARY KEY,
    email                     varchar(255),
    email_confirm_token       varchar(255),
    email_confirmed           boolean,
    username                  varchar(255),
    ip                        varchar(255) NOT NULL,
    authtoken                 varchar(255) NOT NULL,
    created_at                timestamp,
    updated_at                timestamp
);
