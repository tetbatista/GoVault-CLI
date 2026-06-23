package usecase

import (
	"errors"

	"github.com/tetbatista/govault/internal/domain"
)

type SearchCredentialUseCase struct {
	credRepo  domain.CredentialRepository
	userRepo  domain.UserRepository
	deriveKey func(password, salt string) ([]byte, error)
	decryptFn func(ciphertext string, key []byte) (string, error)
}

func NewSearchCredentialUseCase(
	credRepo domain.CredentialRepository,
	userRepo domain.UserRepository,
	deriveKey func(string, string) ([]byte, error),
	decryptFn func(string, []byte) (string, error),
) *SearchCredentialUseCase {
	return &SearchCredentialUseCase{credRepo: credRepo, userRepo: userRepo, deriveKey: deriveKey, decryptFn: decryptFn}
}

type SearchCredentialInput struct {
	UserID     string
	Username   string
	MasterPass string
	Site       string
}

func (uc *SearchCredentialUseCase) Execute(input SearchCredentialInput) (*CredentialOutput, error) {
	user, err := uc.userRepo.FindByUsername(input.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	key, err := uc.deriveKey(input.MasterPass, user.Salt)
	if err != nil {
		return nil, err
	}

	cred, err := uc.credRepo.FindByUserIDAndSite(input.UserID, input.Site)
	if err != nil {
		return nil, err
	}
	if cred == nil {
		return nil, nil
	}

	plainPass, err := uc.decryptFn(cred.Password, key)
	if err != nil {
		return nil, err
	}

	return &CredentialOutput{
		Site:     cred.Site,
		Username: cred.Username,
		Password: plainPass,
	}, nil
}
