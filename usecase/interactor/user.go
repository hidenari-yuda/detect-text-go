package interactor

import (
	"fmt"

	"github.com/hidenari-yuda/umerun-resume/domain/entity"
	"github.com/hidenari-yuda/umerun-resume/usecase"
)

type UserInteractor interface {
	// Gest API
	SignUp(input SignUpInput) (output SignUpOutput, err error)
	SignIn(input SignInInput) (output SignInOutput, err error)
	GetByFirebaseToken(input GetByFirebaseTokenInput) (output GetByFirebaseTokenOutput, err error)
	GetByLineUserId(input GetByLineUserIdInput) (output GetByLineUserIdOutput, err error)

	// line
	GetLineWebHook(input GetLineWebHookInput) (output GetLineWebHookOutput, err error)

	// resume
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

type SignUpInput struct {
	Param *entity.SignUpParam
}

type SignUpOutput struct {
	Ok bool
}

func (i *UserInteractorImpl) SignUp(input SignUpInput) (output SignUpOutput, err error) {
	// ユーザー登録
	err = i.userRepository.SignUp(input.Param)
	if err != nil {
		return output, err
	}

	output.Ok = true

	return output, nil
}

type SignInInput struct {
	Param *entity.SignInParam
}

type SignInOutput struct {
	User *entity.User
}

func (i *UserInteractorImpl) SignIn(input SignInInput) (output SignInOutput, err error) {

	firebaseId, err := i.firebase.VerifyIDToken(input.Param.Token)
	if err != nil {
		return output, err
	}

	fmt.Println("exported firebaseToken is:", input.Param.Token)
	fmt.Println("exported firebaseId is:", firebaseId)

	// ユーザー登録
	output.User, err = i.userRepository.GetByFirebaseId(firebaseId)
	if err != nil {
		err = fmt.Errorf("failed to get user by firebaseId: %w", err)
		return output, err
	}

	return output, nil

}

type GetByFirebaseTokenInput struct {
	Token string
}

type GetByFirebaseTokenOutput struct {
	User *entity.User
}

func (i *UserInteractorImpl) GetByFirebaseToken(input GetByFirebaseTokenInput) (output GetByFirebaseTokenOutput, err error) {

	firebaseId, err := i.firebase.VerifyIDToken(input.Token)
	if err != nil {
		return output, err
	}

	fmt.Println("exported firebaseId is:", firebaseId)

	output.User, err = i.userRepository.GetByFirebaseId(firebaseId)
	if err != nil {
		return output, err
	}

	fmt.Println("exported user is:", output.User)

	return output, nil
}

type GetByLineUserIdInput struct {
	LineUserId string
}

type GetByLineUserIdOutput struct {
	User *entity.User
}

func (i *UserInteractorImpl) GetByLineUserId(input GetByLineUserIdInput) (output GetByLineUserIdOutput, err error) {

	output.User, err = i.userRepository.GetByLineUserId(input.LineUserId)
	if err != nil {
		return output, err
	}

	fmt.Println("exported user is:", output.User)

	return output, nil
}
