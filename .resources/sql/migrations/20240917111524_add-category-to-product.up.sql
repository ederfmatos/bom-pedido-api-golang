ALTER TABLE products ADD COLUMN category_id VARCHAR(36) NOT NULL DEFAULT '';
ALTER TABLE products ALTER COLUMN category_id drop default;
ALTER TABLE products ADD CONSTRAINT fk_product_category FOREIGN KEY (category_id) REFERENCES product_categories;

