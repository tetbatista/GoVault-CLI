package usecase

import (
	"errors"

	"github.com/tetbatista/govault/internal/domain"
)

type LoginUseCase struct {
	userRepo    domain.UserRepository
	checkPassFn func(password, hash string) bool
	createSess  func(userID, username string) error
}

func NewLoginUseCase(
	userRepo domain.UserRepository,
	checkPassFn func(string, string) bool,
	createSess func(string, string) error,
) *LoginUseCase {
	return &LoginUseCase{userRepo: userRepo, checkPassFn: checkPassFn, createSess: createSess}
}

type LoginInput struct {
	Username string
	Password string
}

func (uc *LoginUseCase) Execute(input LoginInput) error {
	user, err := uc.userRepo.FindByUsername(input.Username)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	if !uc.checkPassFn(input.Password, user.Password) {
		return errors.New("incorrect password")
	}

	return uc.createSess(user.ID, user.Username)
}