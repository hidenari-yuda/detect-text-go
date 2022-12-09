-- レシート情報を管理するテーブル
-- +migrate Up
CREATE TABLE IF NOT EXISTS payment_methods (
  id INT NOT NULL AUTO_INCREMENT UNIQUE,
  user_id INT NOT NULL,
  message_id VARCHAR(255) NOT NULL,
  message_type INT NOT NULL, -- 0: text, 1: image, 2: video, 3: audio, 4: file, 5: location, 6: sticker, 7: contact (default: 0)
  text_message VARCHAR(255) NOT NULL, -- メッセージ
  package_id VARCHAR(255) NOT NULL, -- スタンプの表示に使用するID
  sticker_id VARCHAR(255) NOT NULL, -- スタンプの表示に使用するID
  original_content_url VARCHAR(255) NOT NULL, -- 画像ファイル or 動画ファイル or 音声ファイルのUrl
  preview_image_url VARCHAR(255) NOT NULL, -- 画像ファイル or 動画ファイルのプレビュー表示用のファイルUrl
  duration INT, -- 音声ファイルに使用する値
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
