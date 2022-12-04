-- ギフトを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS gifts (
  id INT NOT NULL AUTO_INCREMENT UNIQUE,
  uuid VARCHAR(36) NOT NULL UNIQUE,
  user_id INT NOT NULL,
  receipt_id INT NOT NULL,
  payment_method INT NOT NULL,
  yen INT NOT NULL DEFAULT 0,
  gift_url VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  PRIMARY KEY (id),
  INDEX (user_id),
  INDEX (receipt_id)
);

ALTER TABLE gifts 
ADD CONSTRAINT gifts_user_id_fkey 
FOREIGN KEY (user_id) REFERENCES users(id) 
ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE gifts
ADD CONSTRAINT gifts_receipt_id_fkey
FOREIGN KEY (receipt_id) REFERENCES receipts(id)
ON DELETE CASCADE ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE gifts DROP FOREIGN KEY gifts_user_id_fkey;
ALTER TABLE gifts DROP FOREIGN KEY gifts_receipt_id_fkey;

DROP TABLE IF EXISTS gifts;
