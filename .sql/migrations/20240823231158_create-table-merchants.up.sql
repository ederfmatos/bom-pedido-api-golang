CREATE TABLE merchants
(
    id           VARCHAR(36) PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,
    email        VARCHAR(255) NOT NULL UNIQUE,
    phone_number VARCHAR(11)  NOT NULL UNIQUE,
    tenant_id    VARCHAR(36)  NOT NULL,
    domain       VARCHAR(20)  NOT NULL,
    status       VARCHAR(30)  NOT NULL DEFAULT 'ACTIVE',
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
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
