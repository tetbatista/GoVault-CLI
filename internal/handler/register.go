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

func Register(db *sql.DB) {
	var username string
	fmt.Print("Username: ")
	fmt.Scan(&username)

	fmt.Print("Master password: ")
	passBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		fmt.Println("❌ Error reading password:", err)
		return
	}

	fmt.Print("Confirm master password: ")
	confirmBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		fmt.Println("❌ Error reading confirmation:", err)
		return
	}

	if string(passBytes) != string(confirmBytes) {
		fmt.Println("❌ Passwords do not match.")
		return
	}

	uc := usecase.NewRegisterUseCase(
		repository.NewUserRepository(db),
		auth.HashPassword,
		crypto.GenerateSalt,
	)

	if err := uc.Execute(usecase.RegisterInput{
		Username: username,
		Password: string(passBytes),
	}); err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	fmt.Println("✅ Account created successfully!")
}
