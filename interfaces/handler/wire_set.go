package handler

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewUserHandlerImpl,
	NewPresentHandlerImpl,
)
