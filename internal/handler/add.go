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

func Add(db *sql.DB) {
	auth.RequireAuth(func(claims *auth.Claims) {
		var site, siteUser string
		fmt.Print("Site: ")
		fmt.Scan(&site)
		fmt.Print("Site username: ")
		fmt.Scan(&siteUser)

		fmt.Print("Site password: ")
		sitePassBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			fmt.Println("❌ Error reading site password:", err)
			return
		}

		fmt.Print("Master password (to encrypt): ")
		masterPassBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			fmt.Println("❌ Error reading master password:", err)
			return
		}

		uc := usecase.NewAddCredentialUseCase(
			repository.NewCredentialRepository(db),
			repository.NewUserRepository(db),
			crypto.DeriveKey,
			crypto.Encrypt,
		)

		if err := uc.Execute(usecase.AddCredentialInput{
			UserID:     claims.UserID,
			Username:   claims.Username,
			MasterPass: string(masterPassBytes),
			Site:       site,
			SiteUser:   siteUser,
			SitePass:   string(sitePassBytes),
		}); err != nil {
			fmt.Println("❌ Erro:", err)
			return
		}

		fmt.Println("✅ Credential saved successfully!")
	})
}
