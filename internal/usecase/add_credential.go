package usecase

import (
	"errors"

	"github.com/tetbatista/govault/internal/domain"
)

type AddCredentialUseCase struct {
	credRepo  domain.CredentialRepository
	userRepo  domain.UserRepository
	deriveKey func(password, salt string) ([]byte, error)
	encryptFn func(plaintext string, key []byte) (string, error)
}

func NewAddCredentialUseCase(
	credRepo domain.CredentialRepository,
	userRepo domain.UserRepository,
	deriveKey func(string, string) ([]byte, error),
	encryptFn func(string, []byte) (string, error),
) *AddCredentialUseCase {
	return &AddCredentialUseCase{credRepo: credRepo, userRepo: userRepo, deriveKey: deriveKey, encryptFn: encryptFn}
}

type AddCredentialInput struct {
	UserID     string
	Username   string
	MasterPass string
	Site       string
	SiteUser   string
	SitePass   string
}

func (uc *AddCredentialUseCase) Execute(input AddCredentialInput) error {
	user, err := uc.userRepo.FindByUsername(input.Username)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	key, err := uc.deriveKey(input.MasterPass, user.Salt)
	if err != nil {
		return err
	}

	encryptedPass, err := uc.encryptFn(input.SitePass, key)
	if err != nil {
		return err
	}

	return uc.credRepo.Create(&domain.Credential{
		UserID:   input.UserID,
		Site:     input.Site,
		Username: input.SiteUser,
		Password: encryptedPass,
	})
}
