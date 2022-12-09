-- レシート情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS receipt_pictures (
  id INT NOT NULL AUTO_INCREMENT UNIQUE,
  uuid VARCHAR(36) NOT NULL UNIQUE,
  user_id INT NOT NULL,
  line_user_id VARCHAR(255) NOT NULL,
  url VARCHAR(255) NOT NULL,
  detected_text TEXT NOT NULL,
  service INT,
  payment_service INT,
  total_price INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  PRIMARY KEY (id),
  INDEX (user_id),
  INDEX (line_user_id)
);

ALTER TABLE receipt_pictures 
ADD CONSTRAINT receipt_pictures_user_id_fkey 
FOREIGN KEY (user_id) REFERENCES users(id) 
ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE receipt_pictures
ADD CONSTRAINT receipt_pictures_line_user_id_fkey
FOREIGN KEY (line_user_id) REFERENCES line_users(id)
ON DELETE CASCADE ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE receipt_pictures DROP FOREIGN KEY receipt_pictures_user_id_fkey;
ALTER TABLE receipt_pictures DROP FOREIGN KEY receipt_pictures_line_user_id_fkey;
DROP TABLE IF EXISTS receipt_pictures;
