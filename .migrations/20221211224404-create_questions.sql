-- アンケートを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS questions (
  id INT NOT NULL AUTO_INCREMENT UNIQUE,
  uuid VARCHAR(36) NOT NULL UNIQUE,
  user_id INT NOT NULL,
  line_user_id VARCHAR(255) NOT NULL,
  receipt_picture_id INT NOT NULL,
  question INT NOT NULL,
  text VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  PRIMARY KEY (id),
  INDEX (user_id),
  INDEX (line_user_id)
);

ALTER TABLE questions 
ADD CONSTRAINT questions_user_id_fkey 
FOREIGN KEY (user_id) REFERENCES users(id) 
ON DELETE CASCADE ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE questions DROP FOREIGN KEY questions_user_id_fkey;

DROP TABLE IF EXISTS questions;
