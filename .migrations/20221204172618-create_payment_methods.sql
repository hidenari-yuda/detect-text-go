-- レシート情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS payment_methods (
  id INT NOT NULL AUTO_INCREMENT UNIQUE,
  uuid VARCHAR(36) NOT NULL UNIQUE,
  user_id INT NOT NULL,
  method INT NOT NULL,  -- 支払い方法 0: cash, 1: credit card 2:mobile (default: 0)
  name VARCHAR(255) NOT NULL,
  number VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  PRIMARY KEY (id),
  INDEX (user_id)
);

ALTER TABLE payment_methods 
ADD CONSTRAINT payment_methods_user_id_fkey 
FOREIGN KEY (user_id) REFERENCES users(id) 
ON DELETE CASCADE ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE payment_methods DROP FOREIGN KEY payment_methods_user_id_fkey;
DROP TABLE IF EXISTS payment_methods;
