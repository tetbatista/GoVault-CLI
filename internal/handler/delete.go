package handler

import (
	"database/sql"
	"fmt"

	"github.com/tetbatista/govault/internal/infrastructure/auth"
	"github.com/tetbatista/govault/internal/repository"
	"github.com/tetbatista/govault/internal/usecase"
)

func Delete(db *sql.DB, site string) {
	auth.RequireAuth(func(claims *auth.Claims) {
		var confirm string
		fmt.Printf("Are you sure you want to delete '%s'? (y/n): ", site)
		fmt.Scan(&confirm)

		if confirm != "y" {
			fmt.Println("Operation cancelled.")
			return
		}

		uc := usecase.NewDeleteCredentialUseCase(
			repository.NewCredentialRepository(db),
		)

		if err := uc.Execute(usecase.DeleteCredentialInput{
			UserID: claims.UserID,
			Site:   site,
		}); err != nil {
			fmt.Println("❌ Erro:", err)
			return
		}

		fmt.Println("✅ Credential deleted.")
	})
}