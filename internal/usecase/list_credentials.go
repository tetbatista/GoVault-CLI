package usecase

import (
	"errors"

	"github.com/tetbatista/govault/internal/domain"
)

type ListCredentialsUseCase struct {
	credRepo  domain.CredentialRepository
	userRepo  domain.UserRepository
	deriveKey func(password, salt string) ([]byte, error)
	decryptFn func(ciphertext string, key []byte) (string, error)
}

func NewListCredentialsUseCase(
	credRepo domain.CredentialRepository,
	userRepo domain.UserRepository,
	deriveKey func(string, string) ([]byte, error),
	decryptFn func(string, []byte) (string, error),
) *ListCredentialsUseCase {
	return &ListCredentialsUseCase{credRepo: credRepo, userRepo: userRepo, deriveKey: deriveKey, decryptFn: decryptFn}
}

type ListCredentialsInput struct {
	UserID     string
	Username   string
	MasterPass string
}

type CredentialOutput struct {
	Site     string
	Username string
	Password string
}

func (uc *ListCredentialsUseCase) Execute(input ListCredentialsInput) ([]CredentialOutput, error) {
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

	credentials, err := uc.credRepo.FindAllByUserID(input.UserID)
	if err != nil {
		return nil, err
	}

	var result []CredentialOutput
	for _, cred := range credentials {
		plainPass, err := uc.decryptFn(cred.Password, key)
		if err != nil {
			return nil, err
		}
		result = append(result, CredentialOutput{
			Site:     cred.Site,
			Username: cred.Username,
			Password: plainPass,
		})
	}

	return result, nil
}
