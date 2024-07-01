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

INSERT INTO customers (id, name, email, phone_number, status)
VALUES ('019078bc-cab8-789a-a1e7-4ba2a09561a6', 'Eder Matos', 'ederfmatos@gmail.com', '11999999999', 'ACTIVE');

CREATE TABLE IF NOT EXISTS products
(
    id          VARCHAR(36)                            NOT NULL PRIMARY KEY,
    name        VARCHAR(255)                           NOT NULL,
    description MEDIUMTEXT,
    price       DECIMAL(6, 2)                          NOT NULL,
    status      ENUM ('ACTIVE', 'INACTIVE', 'DELETED') NOT NULL,
    created_at  TIMESTAMP                              NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS orders
(
    id                VARCHAR(36) NOT NULL PRIMARY KEY,
    code              INTEGER     NOT NULL AUTO_INCREMENT UNIQUE,
    customer_id       VARCHAR(36) NOT NULL,
    payment_method    VARCHAR(30) NOT NULL,
    payment_mode      VARCHAR(30) NOT NULL,
    delivery_mode     VARCHAR(30) NOT NULL,
    status            VARCHAR(30) NOT NULL,
    credit_card_token VARCHAR(255),
    `change`          DECIMAL(6, 2),
    delivery_time     TIMESTAMP   NOT NULL,
    created_at        TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_orders_customer FOREIGN KEY (customer_id) REFERENCES customers (id)
);

CREATE TABLE IF NOT EXISTS order_items
(
    id          VARCHAR(36) NOT NULL PRIMARY KEY,
    order_id    VARCHAR(36) NOT NULL,
    product_id  VARCHAR(36) NOT NULL,
    status      VARCHAR(30) NOT NULL,
    quantity    NUMERIC     NOT NULL,
    observation TEXT,
    price       DECIMAL(6, 2),
    created_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_order_items_order FOREIGN KEY (order_id) REFERENCES orders (id),
    CONSTRAINT fk_order_items_product FOREIGN KEY (product_id) REFERENCES products (id)
);
