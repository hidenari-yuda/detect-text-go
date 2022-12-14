-- レシート情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS receipts (
  id INT NOT NULL AUTO_INCREMENT UNIQUE,
  uuid VARCHAR(36) NOT NULL UNIQUE,
  receipt_picture_id INT NOT NULL,
  store_name VARCHAR(255) NOT NULL,
  total_price INT NOT NULL DEFAULT 0,
  purchased_at DATETIME,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  PRIMARY KEY (id),
  INDEX (receipt_picture_id)
);

ALTER TABLE receipts 
ADD CONSTRAINT receipts_receipt_picture_id_fkey
FOREIGN KEY (receipt_picture_id) REFERENCES receipt_pictures(id)
ON DELETE CASCADE ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE receipts DROP FOREIGN KEY receipts_receipt_picture_id_fkey;
DROP TABLE IF EXISTS receipts;
