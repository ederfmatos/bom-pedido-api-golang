CREATE TABLE admins
(
    id          VARCHAR(36) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL UNIQUE,
    merchant_id VARCHAR(36)  NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO admins (id, name, email, merchant_id)
VALUES ('019078bc-cab8-789a-a1e7-4ba2a09561a6', 'Eder Matos', 'ederfmatos@gmail.com', '019078bc-cab8-789a-a1e7-4ba2a09561a6');