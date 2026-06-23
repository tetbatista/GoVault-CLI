package handler

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/tetbatista/govault/internal/infrastructure/auth"
	"github.com/tetbatista/govault/internal/repository"
	"github.com/tetbatista/govault/internal/usecase"
	"golang.org/x/term"
)

func Login(db *sql.DB) {
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

	uc := usecase.NewLoginUseCase(
		repository.NewUserRepository(db),
		auth.CheckPassword,
		auth.CreateSession,
	)

	if err := uc.Execute(usecase.LoginInput{
		Username: username,
		Password: string(passBytes),
	}); err != nil {
		fmt.Println("❌", err)
		return
	}

	fmt.Println("✅ Login successful! Session valid for 30 minutes.")
}