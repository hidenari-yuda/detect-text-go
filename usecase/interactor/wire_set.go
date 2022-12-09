package interactor

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewUserInteractorImpl,
	// NewReceiptInteractorImpl,
	// NewParchasedItemInteractorImpl,
	NewPresentInteractorImpl,
	// NewPaymentMethodInteractorImpl,
	// NewAspInteractorImpl,
	// NewLineMessageInteractorImpl,
)
