-- アンケートを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS questionnaires (
  id INT NOT NULL AUTO_INCREMENT UNIQUE,
  uuid VARCHAR(36) NOT NULL UNIQUE,
  user_id INT NOT NULL,
  line_user_id VARCHAR(255) NOT NULL,
  type INT NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  PRIMARY KEY (id),
  INDEX (user_id),
  INDEX (line_user_id),
  INDEX (receipt_picture_id)
);

ALTER TABLE questionnaires 
ADD CONSTRAINT questionnaires_user_id_fkey 
FOREIGN KEY (user_id) REFERENCES users(id) 
ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE questionnaires
ADD CONSTRAINT questionnaires_receipt_picture_id_fkey
FOREIGN KEY (receipt_picture_id) REFERENCES receipt_pictures(id)
ON DELETE CASCADE ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE questionnaires DROP FOREIGN KEY questionnaires_user_id_fkey;
ALTER TABLE questionnaires DROP FOREIGN KEY questionnaires_receipt_picture_id_fkey;

DROP TABLE IF EXISTS questionnaires;
