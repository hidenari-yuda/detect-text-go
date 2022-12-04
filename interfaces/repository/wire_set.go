package repository

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewUserRepositoryImpl,
	NewReceiptRepositoryImpl,
	NewParchasedItemRepositoryImpl,
	NewGiftRepositoryImpl,
	NewPaymentMethodRepositoryImpl,
	NewAspRepositoryImpl,
	NewLineMessageRepositoryImpl,
)
