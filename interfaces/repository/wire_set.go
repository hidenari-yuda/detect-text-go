package repository

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewUserRepositoryImpl,
	NewLineMessageRepositoryImpl,
	NewPresentRepositoryImpl,
	NewPaymentMethodRepositoryImpl,
	NewCampaignRepositoryImpl,
	NewQuestionRepositoryImpl,
	NewQuestionSelectionRepositoryImpl,

	// レシート関連
	NewReceiptPictureRepositoryImpl,
	NewReceiptRepositoryImpl,
	NewParchasedItemRepositoryImpl,
)
