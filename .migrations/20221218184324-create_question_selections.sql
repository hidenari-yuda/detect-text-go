-- アンケートを管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS question_selections (
  id INT NOT NULL AUTO_INCREMENT UNIQUE,
  uuid VARCHAR(36) NOT NULL UNIQUE,
  question_id INT NOT NULL,
  selection INT NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  PRIMARY KEY (id),
  INDEX (question_id)
);

ALTER TABLE question_selections 
ADD CONSTRAINT question_selections_question_id_fkey
FOREIGN KEY (question_id) REFERENCES questions(id)
ON DELETE CASCADE ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE question_selections DROP FOREIGN KEY question_selections_question_id_fkey;

DROP TABLE IF EXISTS question_selections;
