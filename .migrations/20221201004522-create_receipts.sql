-- レシート情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS receipts (
  id INT NOT NULL AUTO_INCREMENT UNIQUE,
  uuid VARCHAR(36) NOT NULL UNIQUE,
  user_id INT NOT NULL,
  line_user_id VARCHAR(255) NOT NULL,
  image_url VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  PRIMARY KEY (id),
  INDEX (user_id),
  INDEX (line_user_id)
);

ALTER TABLE receipts 
ADD CONSTRAINT receipts_user_id_fkey 
FOREIGN KEY (user_id) REFERENCES users(id) 
ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE receipts
ADD CONSTRAINT receipts_line_user_id_fkey
FOREIGN KEY (line_user_id) REFERENCES line_users(id)
ON DELETE CASCADE ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE receipts DROP FOREIGN KEY receipts_user_id_fkey;
ALTER TABLE receipts DROP FOREIGN KEY receipts_line_user_id_fkey;
DROP TABLE IF EXISTS receipts;
