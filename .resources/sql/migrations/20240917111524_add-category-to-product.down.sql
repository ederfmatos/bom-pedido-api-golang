ALTER TABLE products DROP CONSTRAINT IF EXISTS fk_product_category;
ALTER TABLE products DROP COLUMN IF EXISTS category_id;

