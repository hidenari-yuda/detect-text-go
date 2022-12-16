-- ギフトを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS presents (
  id INT NOT NULL AUTO_INCREMENT UNIQUE,
  uuid VARCHAR(36) NOT NULL UNIQUE,
  user_id INT NOT NULL,
  line_user_id VARCHAR(255) NOT NULL,
  receipt_picture_id INT NOT NULL,
  payment_service INT NOT NULL DEFAULT 0,
  point INT NOT NULL DEFAULT 1,
  url VARCHAR(255) NOT NULL,
  expirary DATETIME NOT NULL,
  used BOOLEAN NOT NULL DEFAULT FALSE,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  PRIMARY KEY (id),
  INDEX (user_id),
  INDEX (receipt_picture_id)
);

ALTER TABLE presents 
ADD CONSTRAINT presents_user_id_fkey 
FOREIGN KEY (user_id) REFERENCES users(id) 
ON DELETE CASCADE ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE presents DROP FOREIGN KEY presents_user_id_fkey;

DROP TABLE IF EXISTS presents;
