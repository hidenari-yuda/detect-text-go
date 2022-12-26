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
  question_progress INT NOT NULL DEFAULT 0, -- 0: 未回答, 1: 回答済み
  prefecture INT NOT NULL DEFAULT 99, -- 0: hokkaido, 1: aomori, 2: iwate, 3: miyagi, 4: akita, 5: yamagata, 6: fukushima, 7: ibaraki, 8: tochigi, 9: gunma, 10: saitama, 11: chiba, 12: tokyo, 13: kanagawa, 14: niigata, 15: toyama, 16: ishikawa, 17: fukui, 18: yamanashi, 19: nagano, 20: gifu, 21: shizuoka, 22: aichi, 23: mie, 24: shiga, 25: kyoto, 26: osaka, 27: hyogo, 28: nara, 29: wakayama, 30: tottori, 31: shimane, 32: okayama, 33: hiroshima, 34: yamaguchi, 35: tokushima, 36: kagawa, 37: ehime, 38: kochi, 39: fukuoka, 40: saga, 41: nagasaki, 42: kumamoto, 43: oita, 44: miyazaki, 45: kagoshima, 46: okinawa
  age INT NOT NULL DEFAULT 99, -- 0: 10代, 1: 20代, 2: 30代, 3: 40代, 4: 50代, 5: 60代, 6: 70代以上
  gender INT NOT NULL DEFAULT 99, -- 0: 男性, 1: 女性 2: その他
  occupation INT NOT NULL DEFAULT 99, -- 0: 学生, 1: 会社員, 2: 自営業, 3: 公務員, 4: その他
  married INT NOT NULL DEFAULT 99, -- 0: 未婚, 1: 既婚
  annual_income INT NOT NULL DEFAULT 99, -- 0: 100万円台, 1: 200万円台, 2: 300万円台, 3: 400万円台, 4: 500万円台, 5: 600万円台, 6: 700万円台, 7: 800万円台, 8: 900万円台, 9: 1000万円以上
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
