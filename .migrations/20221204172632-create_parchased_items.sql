-- ギフトを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS parchased_items (
  id INT NOT NULL AUTO_INCREMENT UNIQUE,
  uuid VARCHAR(36) NOT NULL UNIQUE,
  receipt_id INT NOT NULL,
  name VARCHAR(255) NOT NULL,
  price INT NOT NULL DEFAULT 0,
  number INT NOT NULL DEFAULT 0,
  detected_text VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  PRIMARY KEY (id),
  INDEX (receipt_id)
);

ALTER TABLE parchased_items
ADD CONSTRAINT parchased_items_receipt_id_fkey
FOREIGN KEY (receipt_id) REFERENCES receipts(id)
ON DELETE CASCADE ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE parchased_items DROP FOREIGN KEY parchased_items_receipt_id_fkey;

DROP TABLE IF EXISTS parchased_items;
