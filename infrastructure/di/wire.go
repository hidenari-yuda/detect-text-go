//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/hidenari-yuda/paychan/interfaces"
	"github.com/hidenari-yuda/paychan/interfaces/handler"
	"github.com/hidenari-yuda/paychan/interfaces/repository"
	"github.com/hidenari-yuda/paychan/usecase"
	"github.com/hidenari-yuda/paychan/usecase/interactor"
)

var wireSet = wire.NewSet(
	handler.WireSet,
	interactor.WireSet,
	repository.WireSet,
)

/**
	Handler
**/

// User
//
func InitializeUserHandler(db interfaces.SQLExecuter, fb usecase.Firebase) (h handler.UserHandler) {
	wire.Build(wireSet)
	return
}

// Present
//
func InitializePresentHandler(db interfaces.SQLExecuter, fb usecase.Firebase) (h handler.PresentHandler) {
	wire.Build(wireSet)
	return
}

/**
	Interactor
**/

// User
//
func InitializeUserInteractor(db interfaces.SQLExecuter, fb usecase.Firebase) (i interactor.UserInteractor) {
	wire.Build(wireSet)
	return
}

// Present
//
func InitializePresentInteractor(db interfaces.SQLExecuter, fb usecase.Firebase) (i interactor.PresentInteractor) {
	wire.Build(wireSet)
	return
}
