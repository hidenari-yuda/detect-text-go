package entity

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type LineMessage struct {
	Id                 uint      `db:"id" json:"id"`
	UserId             uint      `db:"user_id" json:"user_id"`
	LineUserId         string    `db:"line_user_id" json:"line_user_id"`
	MessageId          string    `db:"message_id" json:"message_id"`
	MessageType        uint      `db:"message_type" json:"message_type"`                 // 0: text, 1: image, 2: video, 3: audio, 4: file, 5: location, 6: sticker, 7: contact (default: 0)
	TextMessage        string    `db:"text_message" json:"text_message"`                 // メッセージ
	PackageID          string    `db:"package_id" json:"package_id"`                     // スタンプの表示に使用するID
	StickerID          string    `db:"sticker_id" json:"sticker_id"`                     // スタンプの表示に使用するID
	OriginalContentUrl string    `db:"original_content_url" json:"original_content_url"` // 画像ファイル or 動画ファイル or 音声ファイルのUrl
	PreviewImageUrl    string    `db:"preview_image_url" json:"preview_image_url"`       // 画像ファイル or 動画ファイルのプレビュー表示用のファイルUrl
	Duration           null.Int  `db:"duration" json:"duration"`                         // 音声ファイルに使用する値
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
}
