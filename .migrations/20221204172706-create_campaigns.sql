-- 広告情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS campaigns (
  id INT NOT NULL AUTO_INCREMENT UNIQUE,
  uuid VARCHAR(36) NOT NULL UNIQUE,
  service INT NOT NULL,
  url VARCHAR(255) NOT NULL,
  picture_url VARCHAR(2000) NOT NULL,
  price INT NOT NULL DEFAULT 0,
  title VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  impression INT NOT NULL DEFAULT 0,
  click INT NOT NULL DEFAULT 0,
  client_id INT NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  PRIMARY KEY (id)
);

-- +migrate Down
ALTER TABLE ads DROP FOREIGN KEY ads_user_id_fkey;
DROP TABLE IF EXISTS ads;