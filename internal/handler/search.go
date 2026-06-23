package handler

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/tetbatista/govault/internal/infrastructure/auth"
	"github.com/tetbatista/govault/internal/infrastructure/crypto"
	"github.com/tetbatista/govault/internal/repository"
	"github.com/tetbatista/govault/internal/usecase"
	"golang.org/x/term"
)

func Search(db *sql.DB, site string) {
	auth.RequireAuth(func(claims *auth.Claims) {
		fmt.Print("Master password (to decrypt): ")
		masterPassBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			fmt.Println("❌ Error reading master password:", err)
			return
		}

		uc := usecase.NewSearchCredentialUseCase(
			repository.NewCredentialRepository(db),
			repository.NewUserRepository(db),
			crypto.DeriveKey,
			crypto.Decrypt,
		)

		result, err := uc.Execute(usecase.SearchCredentialInput{
			UserID:     claims.UserID,
			Username:   claims.Username,
			MasterPass: string(masterPassBytes),
			Site:       site,
		})
		if err != nil {
			fmt.Println("❌ Erro:", err)
			return
		}
		if result == nil {
			fmt.Println("❌ No credential found for:", site)
			return
		}

		fmt.Printf("\n%-18s %-25s %s\n", "Site", "Username", "Password")
		fmt.Println("--------------------------------------------------------------")
		fmt.Printf("%-18s %-25s %s\n\n", result.Site, result.Username, result.Password)
	})
}