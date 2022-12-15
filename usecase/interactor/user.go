package interactor

import (
	"fmt"

	"github.com/hidenari-yuda/paychan-server/domain/entity"
	"github.com/hidenari-yuda/paychan-server/usecase"
)

type UserInteractor interface {
	// Gest API
	SignUp(param *entity.SignUpParam) (ok bool, err error)
	SignIn(param *entity.SignInParam) (user *entity.User, err error)
	GetByFirebaseToken(token string) (user *entity.User, err error)
	GetByLineUserId(lineUserId string) (user *entity.User, err error)

	// admin
	GetAll() (users []*entity.User, err error)

	// line
	GetLineWebHook(param *entity.LineWebHook) (ok bool, err error)
}

type UserInteractorImpl struct {
	firebase                 usecase.Firebase
	userRepository           usecase.UserRepository
	receiptPictureRepository usecase.ReceiptPictureRepository
	receiptRepository        usecase.ReceiptRepository
	parchasedItemRepository  usecase.ParchasedItemRepository
	paymentMethodRepository  usecase.PaymentMethodRepository
	presentRepository        usecase.PresentRepository
	lineMessageRepository    usecase.LineMessageRepository
	aspRepository            usecase.AspRepository
}

func NewUserInteractorImpl(
	fb usecase.Firebase,
	uR usecase.UserRepository,
	rpR usecase.ReceiptPictureRepository,
	rR usecase.ReceiptRepository,
	piR usecase.ParchasedItemRepository,
	pmR usecase.PaymentMethodRepository,
	pR usecase.PresentRepository,
	lmR usecase.LineMessageRepository,
	aR usecase.AspRepository,
) UserInteractor {
	return &UserInteractorImpl{
		firebase:                 fb,
		userRepository:           uR,
		receiptPictureRepository: rpR,
		receiptRepository:        rR,
		parchasedItemRepository:  piR,
		paymentMethodRepository:  pmR,
		presentRepository:        pR,
		lineMessageRepository:    lmR,
		aspRepository:            aR,
	}
}

func (i *UserInteractorImpl) SignUp(param *entity.SignUpParam) (ok bool, err error) {
	// ユーザー登録
	err = i.userRepository.SignUp(param)
	if err != nil {
		return ok, err
	}

	ok = true

	return ok, nil
}

func (i *UserInteractorImpl) SignIn(param *entity.SignInParam) (user *entity.User, err error) {

	firebaseId, err := i.firebase.VerifyIDToken(param.Token)
	if err != nil {
		return user, err
	}

	fmt.Println("exported firebaseToken is:", param.Token)
	fmt.Println("exported firebaseId is:", firebaseId)

	// ユーザー登録
	user, err = i.userRepository.GetByFirebaseId(firebaseId)
	if err != nil {
		err = fmt.Errorf("failed to get user by firebaseId: %w", err)
		return user, err
	}

	return user, nil

}

func (i *UserInteractorImpl) GetByFirebaseToken(token string) (user *entity.User, err error) {

	firebaseId, err := i.firebase.VerifyIDToken(token)
	if err != nil {
		return user, err
	}

	fmt.Println("exported firebaseId is:", firebaseId)

	user, err = i.userRepository.GetByFirebaseId(firebaseId)
	if err != nil {
		return user, err
	}

	fmt.Println("exported user is:", user)

	return user, nil
}

func (i *UserInteractorImpl) GetAll() (users []*entity.User, err error) {

	users, err = i.userRepository.GetAll()
	if err != nil {
		return users, err
	}

	return users, nil
}

func (i *UserInteractorImpl) GetByLineUserId(lineUserId string) (user *entity.User, err error) {

	user, err = i.userRepository.GetByLineUserId(lineUserId)
	if err != nil {
		return user, err
	}

	fmt.Println("exported user is:", user)

	return user, nil
}
