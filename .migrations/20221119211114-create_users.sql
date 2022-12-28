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
  industry INT NOT NULL DEFAULT 99, -- 0: 農林水産業, 1: 鉱業, 2: 建設業, 3: 製造業, 4: 電気・ガス・熱供給・水道業, 5: 情報通信業, 6: 小売業, 7: 飲食業, 8: 金融業, 9: 不動産業, 10: サービス業, 11: 教育・学習支援業, 12: 医療・福祉, 13: 記録媒体の製作・配給, 14: 観光・宿泊業, 15: その他
  living_with INT NOT NULL DEFAULT 99, -- 0: 一人暮らし, 1: 二人暮らし, 2: 三人暮らし, 3: その他
  marriage INT NOT NULL DEFAULT 99, -- 0: 未婚, 1: 既婚
  child INT NOT NULL DEFAULT 99, -- 0: 子供なし, 1: 子供あり
  child_age INT NOT NULL DEFAULT 99, -- 0: 0歳, 1: 1歳, 2: 2歳, 3: 3歳, 4: 4歳, 5: 5歳, 6: 6歳, 7: 7歳, 8: 8歳, 9: 9歳, 10: 10歳, 11: 11歳, 12: 12歳, 13: 13歳, 14: 14歳, 15: 15歳, 16: 16歳, 17: 17歳, 18: 18歳以上
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
