package usecase

import "github.com/tetbatista/govault/internal/domain"

type DeleteCredentialUseCase struct {
	credRepo domain.CredentialRepository
}

func NewDeleteCredentialUseCase(credRepo domain.CredentialRepository) *DeleteCredentialUseCase {
	return &DeleteCredentialUseCase{credRepo: credRepo}
}

type DeleteCredentialInput struct {
	UserID string
	Site   string
}

func (uc *DeleteCredentialUseCase) Execute(input DeleteCredentialInput) error {
	return uc.credRepo.DeleteByUserIDAndSite(input.UserID, input.Site)
}
