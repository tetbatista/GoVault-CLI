package usecase

import (
	"errors"

	"github.com/tetbatista/govault/internal/domain"
)

type RegisterUseCase struct {
	userRepo domain.UserRepository
	hashFn   func(string) (string, error)
	saltFn   func() (string, error)
}

func NewRegisterUseCase(
	userRepo domain.UserRepository,
	hashFn func(string) (string, error),
	saltFn func() (string, error),
) *RegisterUseCase {
	return &RegisterUseCase{userRepo: userRepo, hashFn: hashFn, saltFn: saltFn}
}

type RegisterInput struct {
	Username string
	Password string
}

func (uc *RegisterUseCase) Execute(input RegisterInput) error {
	existing, err := uc.userRepo.FindByUsername(input.Username)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("user already exists")
	}

	salt, err := uc.saltFn()
	if err != nil {
		return err
	}

	hash, err := uc.hashFn(input.Password)
	if err != nil {
		return err
	}

	return uc.userRepo.Create(&domain.User{
		Username: input.Username,
		Password: hash,
		Salt:     salt,
	})
}