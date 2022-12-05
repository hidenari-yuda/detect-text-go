package repository

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewUserRepositoryImpl,
	NewLineMessageRepositoryImpl,
	NewGiftRepositoryImpl,
	NewPaymentMethodRepositoryImpl,
	NewAspRepositoryImpl,

	// レシート関連
	NewReceiptPictureRepositoryImpl,
	NewReceiptRepositoryImpl,
	NewParchasedItemRepositoryImpl,
)
