CREATE DATABASE IF NOT EXISTS bompedido;
CREATE SCHEMA IF NOT EXISTS bompedido;
USE bompedido;

CREATE TABLE IF NOT EXISTS customers
(
    id           VARCHAR(36)                NOT NULL PRIMARY KEY,
    name         VARCHAR(255)               NOT NULL,
    email        VARCHAR(255)               NOT NULL UNIQUE,
    phone_number VARCHAR(11)                NOT NULL UNIQUE,
    status       ENUM ('ACTIVE', 'DELETED') NOT NULL,
    created_at   TIMESTAMP                  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS products
(
    id          VARCHAR(36)                            NOT NULL PRIMARY KEY,
    name        VARCHAR(255)                           NOT NULL,
    description MEDIUMTEXT,
    price       DECIMAL(6, 2)                          NOT NULL,
    status      ENUM ('ACTIVE', 'INACTIVE', 'DELETED') NOT NULL,
    created_at  TIMESTAMP                              NOT NULL DEFAULT CURRENT_TIMESTAMP
);