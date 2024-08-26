CREATE TABLE merchants
(
    id           VARCHAR(36) PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,
    email        VARCHAR(255) NOT NULL,
    phone_number VARCHAR(11)  NOT NULL,
    tenant_id    VARCHAR(36)  NOT NULL UNIQUE,
    domain       VARCHAR(20)  NOT NULL,
    status       VARCHAR(30)  NOT NULL DEFAULT 'ACTIVE',
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_merchant_email UNIQUE (email, tenant_id),
    CONSTRAINT uk_merchant_phone_number UNIQUE (phone_number, tenant_id)
);
CREATE INDEX ids_merchants_domain ON merchants (domain);
CREATE INDEX ids_merchants_tenant_id ON merchants (tenant_id);

CREATE TABLE merchant_address
(
    id           SERIAL PRIMARY KEY,
    merchant_id  VARCHAR(36)  NOT NULL,
    street       VARCHAR(255) NOT NULL,
    number       VARCHAR(20),
    neighborhood VARCHAR(255),
    postal_code  VARCHAR(8),
    city         VARCHAR(100) NOT NULL,
    state        VARCHAR(2)   NOT NULL,
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_merchant_address_merchant FOREIGN KEY (merchant_id) REFERENCES merchants (id)
);

CREATE INDEX ids_merchant_address_merchant_id ON merchant_address (merchant_id);

CREATE TABLE merchant_opening_hour
(
    id           SERIAL PRIMARY KEY,
    merchant_id  VARCHAR(36) NOT NULL,
    day_of_week  NUMERIC(2)  NOT NULL,
    initial_time TIME        NOT NULL,
    final_time   TIME        NOT NULL,
    created_at   TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_merchant_opening_hour_merchant FOREIGN KEY (merchant_id) REFERENCES merchants (id)
);
CREATE INDEX ids_merchant_opening_hour_merchant_id ON merchant_opening_hour (merchant_id);

CREATE TABLE merchant_payment_gateway_configs
(
    id          SERIAL      NOT NULL PRIMARY KEY,
    merchant_id VARCHAR(36) NOT NULL,
    gateway     VARCHAR(50) NOT NULL,
    credentials TEXT        NOT NULL,
    CONSTRAINT fk_merchant_payment_gateway_configs_merchant FOREIGN KEY (merchant_id) REFERENCES merchants (id)
);
CREATE INDEX ids_merchant_payment_gateway_configs_merchant_id ON merchant_payment_gateway_configs (merchant_id);

CREATE TABLE admins
(
    id          VARCHAR(36) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL,
    merchant_id VARCHAR(36)  NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_admins_merchant FOREIGN KEY (merchant_id) REFERENCES merchants (id),
    CONSTRAINT uk_admin_email UNIQUE (email, merchant_id)
);
CREATE INDEX ids_admin_email ON admins (email);

CREATE TABLE customers
(
    id           VARCHAR(36)  NOT NULL PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,
    email        VARCHAR(255) NOT NULL,
    phone_number VARCHAR(11),
    status       VARCHAR(20)  NOT NULL,
    tenant_id    VARCHAR(30)  NOT NULL,
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_customer_email UNIQUE (email, tenant_id),
    CONSTRAINT uk_customer_phone_number UNIQUE (phone_number, tenant_id)
);
CREATE INDEX ids_customer_email_tenant ON customers (email, tenant_id);

CREATE TABLE products
(
    id          VARCHAR(36)   NOT NULL PRIMARY KEY,
    name        VARCHAR(255)  NOT NULL,
    description TEXT,
    price       DECIMAL(6, 2) NOT NULL,
    status      VARCHAR(20)   NOT NULL,
    tenant_id   VARCHAR(30)   NOT NULL,
    created_at  TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_product_name UNIQUE (name, tenant_id)
);
CREATE INDEX ids_product_tenant ON products (tenant_id);

CREATE TABLE orders
(
    id                VARCHAR(36)   NOT NULL PRIMARY KEY,
    code              SERIAL        NOT NULL,
    customer_id       VARCHAR(36)   NOT NULL,
    payment_method    VARCHAR(30)   NOT NULL,
    payment_mode      VARCHAR(30)   NOT NULL,
    delivery_mode     VARCHAR(30)   NOT NULL,
    status            VARCHAR(30)   NOT NULL,
    credit_card_token VARCHAR(255),
    payback           DECIMAL(6, 2),
    amount            DECIMAL(6, 2) NOT NULL,
    delivery_time     TIMESTAMP     NOT NULL,
    merchant_id       VARCHAR(36)   NOT NULL,
    created_at        TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_orders_customer FOREIGN KEY (customer_id) REFERENCES customers (id),
    CONSTRAINT fk_orders_merchant FOREIGN KEY (merchant_id) REFERENCES merchants (id)
);
CREATE INDEX ids_order_customer ON orders (customer_id, merchant_id);

CREATE TABLE order_items
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

CREATE TABLE order_history
(
    id         SERIAL      NOT NULL PRIMARY KEY,
    order_id   VARCHAR(36) NOT NULL,
    changed_by VARCHAR(36) NOT NULL,
    changed_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status     VARCHAR(30) NOT NULL,
    data       TEXT,
    CONSTRAINT fk_order_history_order FOREIGN KEY (order_id) REFERENCES orders (id)
);

CREATE TYPE transaction_type AS ENUM ('PIX', 'CREDIT_CARD');

CREATE TABLE transactions
(
    id         VARCHAR(36) PRIMARY KEY NOT NULL,
    order_id   VARCHAR(36)             NOT NULL,
    amount     NUMERIC(6, 2)           NOT NULL,
    status     VARCHAR(20)             NOT NULL,
    created_at TIMESTAMP               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    type       transaction_type        NOT NULL,
    CONSTRAINT fk_transactions_order_id FOREIGN KEY (order_id) REFERENCES orders (id)
);
CREATE INDEX idx_transactions_order ON transactions (order_id);

CREATE TABLE pix_transactions
(
    qr_code         TEXT        NOT NULL,
    qr_code_link    TEXT        NOT NULL,
    payment_gateway VARCHAR(50) NOT NULL
) inherits (transactions);

INSERT INTO merchants (id, name, email, phone_number, tenant_id, domain)
VALUES ('019078bc-cab8-789a-a1e7-4ba2a09561a6', 'Minha Lanchonete', 'ederfmatos@gmail.com', '11999999999',
        '01G65Z755AFWAKHE12NY0CQ9FH', 'minhalanchonete');

INSERT INTO merchant_address (merchant_id, street, number, neighborhood, postal_code, city, state)
VALUES ('019078bc-cab8-789a-a1e7-4ba2a09561a6', 'Avenida Nove de Julho', '100', 'Centro', '00000000', 'São Paulo',
        'SP');

INSERT INTO merchant_opening_hour (merchant_id, day_of_week, initial_time, final_time)
VALUES ('019078bc-cab8-789a-a1e7-4ba2a09561a6', 1, '18:30', '23:30'),
       ('019078bc-cab8-789a-a1e7-4ba2a09561a6', 3, '18:30', '23:30'),
       ('019078bc-cab8-789a-a1e7-4ba2a09561a6', 4, '18:30', '23:30'),
       ('019078bc-cab8-789a-a1e7-4ba2a09561a6', 5, '18:30', '23:30'),
       ('019078bc-cab8-789a-a1e7-4ba2a09561a6', 6, '18:30', '00:00'),
       ('019078bc-cab8-789a-a1e7-4ba2a09561a6', 7, '18:30', '00:30');

INSERT INTO merchant_payment_gateway_configs(merchant_id, gateway, credentials)
VALUES ('019078bc-cab8-789a-a1e7-4ba2a09561a6', 'MERCADO_PAGO', '123456');

INSERT INTO admins (id, name, email, merchant_id)
VALUES ('019078bc-cab8-789a-a1e7-4ba2a09561a6', 'Eder Matos', 'ederfmatos@gmail.com',
        '019078bc-cab8-789a-a1e7-4ba2a09561a6');

INSERT INTO customers (id, name, email, phone_number, status, tenant_id)
VALUES ('019078bc-cab8-789a-a1e7-4ba2a09561a6', 'Eder Matos', 'ederfmatos@gmail.com', '11999999999', 'ACTIVE',
        'minhalanchonete');