//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/hidenari-yuda/detect-text/interfaces"
	"github.com/hidenari-yuda/detect-text/interfaces/handler"
	"github.com/hidenari-yuda/detect-text/interfaces/repository"
	"github.com/hidenari-yuda/detect-text/usecase"
	"github.com/hidenari-yuda/detect-text/usecase/interactor"
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

/**
	Interactor
**/

// User
//
func InitializeUserInteractor(db interfaces.SQLExecuter, fb usecase.Firebase) (i interactor.UserInteractor) {
	wire.Build(wireSet)
	return
}
