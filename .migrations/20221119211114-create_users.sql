-- ユーザー情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
  id INT NOT NULL AUTO_INCREMENT UNIQUE,
  uuid VARCHAR(36) NOT NULL UNIQUE,
  firebase_id VARCHAR(255) NOT NULL UNIQUE,
  line_user_id VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  point INT NOT NULL DEFAULT 0,
  line_name VARCHAR(255) NOT NULL,
  picture_url VARCHAR(255) NOT NULL,
  status_message VARCHAR(255) NOT NULL,
  language VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  PRIMARY KEY (id),
  INDEX (line_user_id)
);

-- +migrate Down
DROP TABLE IF EXISTS users;
