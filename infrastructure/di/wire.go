//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/hidenari-yuda/paychan-server/domain/config"
	"github.com/hidenari-yuda/paychan-server/interfaces"
	"github.com/hidenari-yuda/paychan-server/interfaces/handler"
	"github.com/hidenari-yuda/paychan-server/interfaces/repository"
	"github.com/hidenari-yuda/paychan-server/usecase"
	"github.com/hidenari-yuda/paychan-server/usecase/interactor"
)

var wireSet = wire.NewSet(
	handler.WireSet,
	interactor.WireSet,
	repository.WireSet,
)

/**
	Handler
**/

//Admin
//
func InitializeAdminHandler(db interfaces.SQLExecuter, appConfig config.App) (h handler.AdminHandler) {
	wire.Build(wireSet)
	return
}

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

//Admin
//
func InitializeAdminInteractor(db interfaces.SQLExecuter, appConfig config.App) (i interactor.AdminInteractor) {
	wire.Build(wireSet)
	return
}

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
