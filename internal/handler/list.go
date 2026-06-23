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

func List(db *sql.DB) {
	auth.RequireAuth(func(claims *auth.Claims) {
		fmt.Print("Master password (to decrypt): ")
		masterPassBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			fmt.Println("❌ Error reading master password:", err)
			return
		}

		uc := usecase.NewListCredentialsUseCase(
			repository.NewCredentialRepository(db),
			repository.NewUserRepository(db),
			crypto.DeriveKey,
			crypto.Decrypt,
		)

		credentials, err := uc.Execute(usecase.ListCredentialsInput{
			UserID:     claims.UserID,
			Username:   claims.Username,
			MasterPass: string(masterPassBytes),
		})
		if err != nil {
			fmt.Println("❌ Erro:", err)
			return
		}

		if len(credentials) == 0 {
			fmt.Println("No credentials saved yet.")
			return
		}

		fmt.Printf("\n%-18s %-25s %s\n", "Site", "Username", "Password")
		fmt.Println("--------------------------------------------------------------")
		for _, c := range credentials {
			fmt.Printf("%-18s %-25s %s\n", c.Site, c.Username, c.Password)
		}
		fmt.Println()
	})
}
